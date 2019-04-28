package kazoo

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/samuel/go-zookeeper/zk"
)

// Topic interacts with Kafka's topic metadata in Zookeeper.
type Topic struct {
	Name string
	kz   *Kazoo
}

// TopicList is a type that implements the sortable interface for a list of Topic instances.
type TopicList []*Topic

// Partition interacts with Kafka's partition metadata in Zookeeper.
type Partition struct {
	topic    *Topic
	ID       int32
	Replicas []int32
}

// PartitionList is a type that implements the sortable interface for a list of Partition instances
type PartitionList []*Partition

// Topics returns a list of all registered Kafka topics.
func (kz *Kazoo) Topics() (TopicList, error) {
	root := fmt.Sprintf("%s/brokers/topics", kz.conf.Chroot)
	children, _, err := kz.conn.Children(root)
	if err != nil {
		return nil, err
	}

	result := make(TopicList, 0, len(children))
	for _, name := range children {
		result = append(result, kz.Topic(name))
	}
	return result, nil
}

// WatchTopics returns a list of all registered Kafka topics, and
// watches that list for changes.
func (kz *Kazoo) WatchTopics() (TopicList, <-chan zk.Event, error) {
	root := fmt.Sprintf("%s/brokers/topics", kz.conf.Chroot)
	children, _, c, err := kz.conn.ChildrenW(root)
	if err != nil {
		return nil, nil, err
	}

	result := make(TopicList, 0, len(children))
	for _, name := range children {
		result = append(result, kz.Topic(name))
	}
	return result, c, nil
}

// Topic returns a Topic instance for a given topic name
func (kz *Kazoo) Topic(topic string) *Topic {
	return &Topic{Name: topic, kz: kz}
}

// Exists returns true if the topic exists on the Kafka cluster.
func (t *Topic) Exists() (bool, error) {
	return t.kz.exists(fmt.Sprintf("%s/brokers/topics/%s", t.kz.conf.Chroot, t.Name))
}

// Partitions returns a list of all partitions for the topic.
func (t *Topic) Partitions() (PartitionList, error) {
	node := fmt.Sprintf("%s/brokers/topics/%s", t.kz.conf.Chroot, t.Name)
	value, _, err := t.kz.conn.Get(node)
	if err != nil {
		return nil, err
	}

	return t.parsePartitions(value)
}

// WatchPartitions returns a list of all partitions for the topic, and watches the topic for changes.
func (t *Topic) WatchPartitions() (PartitionList, <-chan zk.Event, error) {
	node := fmt.Sprintf("%s/brokers/topics/%s", t.kz.conf.Chroot, t.Name)
	value, _, c, err := t.kz.conn.GetW(node)
	if err != nil {
		return nil, nil, err
	}

	list, err := t.parsePartitions(value)
	return list, c, err
}

// parsePartitions pases the JSON representation of the partitions
// that is stored as data on the topic node in Zookeeper.
func (t *Topic) parsePartitions(value []byte) (PartitionList, error) {
	type topicMetadata struct {
		Partitions map[string][]int32 `json:"partitions"`
	}

	var tm topicMetadata
	if err := json.Unmarshal(value, &tm); err != nil {
		return nil, err
	}

	result := make(PartitionList, len(tm.Partitions))
	for partitionNumber, replicas := range tm.Partitions {
		partitionID, err := strconv.ParseInt(partitionNumber, 10, 32)
		if err != nil {
			return nil, err
		}

		replicaIDs := make([]int32, 0, len(replicas))
		for _, r := range replicas {
			replicaIDs = append(replicaIDs, int32(r))
		}
		result[partitionID] = t.Partition(int32(partitionID), replicaIDs)
	}

	return result, nil
}

// Partition returns a Partition instance for the topic.
func (t *Topic) Partition(id int32, replicas []int32) *Partition {
	return &Partition{ID: id, Replicas: replicas, topic: t}
}

// Config returns topic-level configuration settings as a map.
func (t *Topic) Config() (map[string]string, error) {
	value, _, err := t.kz.conn.Get(fmt.Sprintf("%s/config/topics/%s", t.kz.conf.Chroot, t.Name))
	if err != nil {
		return nil, err
	}

	var topicConfig struct {
		ConfigMap map[string]string `json:"config"`
	}

	if err := json.Unmarshal(value, &topicConfig); err != nil {
		return nil, err
	}

	return topicConfig.ConfigMap, nil
}

// Topic returns the Topic of this partition.
func (p *Partition) Topic() *Topic {
	return p.topic
}

// Key returns a unique identifier for the partition, using the form "topic/partition".
func (p *Partition) Key() string {
	return fmt.Sprintf("%s/%d", p.topic.Name, p.ID)
}

// PreferredReplica returns the preferred replica for this partition.
func (p *Partition) PreferredReplica() int32 {
	if len(p.Replicas) > 0 {
		return p.Replicas[0]
	} else {
		return -1
	}
}

// Leader returns the broker ID of the broker that is currently the leader for the partition.
func (p *Partition) Leader() (int32, error) {
	if state, err := p.state(); err != nil {
		return -1, err
	} else {
		return state.Leader, nil
	}
}

// ISR returns the broker IDs of the current in-sync replica set for the partition
func (p *Partition) ISR() ([]int32, error) {
	if state, err := p.state(); err != nil {
		return nil, err
	} else {
		return state.ISR, nil
	}
}

func (p *Partition) UnderReplicated() (bool, error) {
	if state, err := p.state(); err != nil {
		return false, err
	} else {
		return len(state.ISR) < len(p.Replicas), nil
	}
}

func (p *Partition) UsesPreferredReplica() (bool, error) {
	if state, err := p.state(); err != nil {
		return false, err
	} else {
		return len(state.ISR) > 0 && state.ISR[0] == p.Replicas[0], nil
	}
}

// partitionState represents the partition state as it is stored as JSON
// in Zookeeper on the partition's state node.
type partitionState struct {
	Leader int32   `json:"leader"`
	ISR    []int32 `json:"isr"`
}

// state retrieves and parses the partition State
func (p *Partition) state() (partitionState, error) {
	var state partitionState
	node := fmt.Sprintf("%s/brokers/topics/%s/partitions/%d/state", p.topic.kz.conf.Chroot, p.topic.Name, p.ID)
	value, _, err := p.topic.kz.conn.Get(node)
	if err != nil {
		return state, err
	}

	if err := json.Unmarshal(value, &state); err != nil {
		return state, err
	}

	return state, nil
}

// Find returns the topic with the given name if it exists in the topic list,
// and will return `nil` otherwise.
func (tl TopicList) Find(name string) *Topic {
	for _, topic := range tl {
		if topic.Name == name {
			return topic
		}
	}
	return nil
}

func (tl TopicList) Len() int {
	return len(tl)
}

func (tl TopicList) Less(i, j int) bool {
	return tl[i].Name < tl[j].Name
}

func (tl TopicList) Swap(i, j int) {
	tl[i], tl[j] = tl[j], tl[i]
}

func (pl PartitionList) Len() int {
	return len(pl)
}

func (pl PartitionList) Less(i, j int) bool {
	return pl[i].topic.Name < pl[j].topic.Name || (pl[i].topic.Name == pl[j].topic.Name && pl[i].ID < pl[j].ID)
}

func (pl PartitionList) Swap(i, j int) {
	pl[i], pl[j] = pl[j], pl[i]
}

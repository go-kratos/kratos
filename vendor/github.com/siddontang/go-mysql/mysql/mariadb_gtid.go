package mysql

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/juju/errors"
)

type MariadbGTID struct {
	DomainID       uint32
	ServerID       uint32
	SequenceNumber uint64
}

// We don't support multi source replication, so the mariadb gtid set may have only domain-server-sequence
func ParseMariadbGTIDSet(str string) (GTIDSet, error) {
	if len(str) == 0 {
		return &MariadbGTID{0, 0, 0}, nil
	}

	seps := strings.Split(str, "-")

	gtid := new(MariadbGTID)

	if len(seps) != 3 {
		return gtid, errors.Errorf("invalid Mariadb GTID %v, must domain-server-sequence", str)
	}

	domainID, err := strconv.ParseUint(seps[0], 10, 32)
	if err != nil {
		return gtid, errors.Errorf("invalid MariaDB GTID Domain ID (%v): %v", seps[0], err)
	}

	serverID, err := strconv.ParseUint(seps[1], 10, 32)
	if err != nil {
		return gtid, errors.Errorf("invalid MariaDB GTID Server ID (%v): %v", seps[1], err)
	}

	sequenceID, err := strconv.ParseUint(seps[2], 10, 64)
	if err != nil {
		return gtid, errors.Errorf("invalid MariaDB GTID Sequence number (%v): %v", seps[2], err)
	}

	return &MariadbGTID{
		DomainID:       uint32(domainID),
		ServerID:       uint32(serverID),
		SequenceNumber: sequenceID}, nil
}

func (gtid *MariadbGTID) String() string {
	if gtid.DomainID == 0 && gtid.ServerID == 0 && gtid.SequenceNumber == 0 {
		return ""
	}

	return fmt.Sprintf("%d-%d-%d", gtid.DomainID, gtid.ServerID, gtid.SequenceNumber)
}

func (gtid *MariadbGTID) Encode() []byte {
	return []byte(gtid.String())
}

func (gtid *MariadbGTID) Equal(o GTIDSet) bool {
	other, ok := o.(*MariadbGTID)
	if !ok {
		return false
	}

	return *gtid == *other
}

func (gtid *MariadbGTID) Contain(o GTIDSet) bool {
	other, ok := o.(*MariadbGTID)
	if !ok {
		return false
	}

	return gtid.DomainID == other.DomainID && gtid.SequenceNumber >= other.SequenceNumber
}

func (gtid *MariadbGTID) Update(GTIDStr string) error {
	newGTID, err := ParseMariadbGTIDSet(GTIDStr)
	if err != nil {
		return err
	}

	*gtid = *(newGTID.(*MariadbGTID))

	return nil
}

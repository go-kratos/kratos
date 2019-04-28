package dao

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"io"
	"os"
	"strings"

	"go-common/library/ecode"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
)

// authConfig get Auth Config.
func (d *Dao) authConfig() (authStr string, err error) {
	authConfig := types.AuthConfig{
		Username:      d.c.BiliHub.Username,
		Password:      d.c.BiliHub.Password,
		ServerAddress: d.c.BiliHub.HostName,
	}

	var encodedJSON []byte
	if encodedJSON, err = json.Marshal(authConfig); err != nil {
		return
	}
	authStr = base64.URLEncoding.EncodeToString(encodedJSON)
	return

}

// ImagePull Image Pull.
func (d *Dao) ImagePull(imageName string) (err error) {
	var (
		out     io.ReadCloser
		authStr string
	)

	if authStr, err = d.authConfig(); err != nil {
		return
	}

	if out, err = d.dockerClient.ImagePull(context.TODO(), imageName, types.ImagePullOptions{RegistryAuth: authStr}); err != nil {
		err = ecode.MerlinImagePullErr
		return
	}

	defer out.Close()
	io.Copy(os.Stdout, out)

	return
}

// ImagePush Image Push.
func (d *Dao) ImagePush(imageName string) (err error) {
	var (
		out     io.ReadCloser
		authStr string
	)

	if authStr, err = d.authConfig(); err != nil {
		return
	}

	if out, err = d.dockerClient.ImagePush(context.TODO(), imageName, types.ImagePushOptions{RegistryAuth: authStr}); err != nil {
		err = ecode.MerlinImagePushErr
		return
	}

	defer out.Close()
	io.Copy(os.Stdout, out)

	return
}

// ImageTag Image Tag.
func (d *Dao) ImageTag(imageSrcName, imageTagName string) (err error) {
	err = d.dockerClient.ImageTag(context.TODO(), imageSrcName, imageTagName)
	return
}

// ImageRemove Image Remove.
func (d *Dao) ImageRemove(imageID string) (err error) {
	_, err = d.dockerClient.ImageRemove(context.TODO(), imageID, types.ImageRemoveOptions{Force: true, PruneChildren: true})
	return
}

// ImageGetID Image Get ID.
func (d *Dao) ImageGetID(imageName string) (imageID string, err error) {
	var images []types.ImageSummary

	if images, err = d.dockerClient.ImageList(context.TODO(), types.ImageListOptions{}); err != nil {
		return
	}

	for _, image := range images {
		for _, tag := range image.RepoTags {
			if tag == imageName {
				imageID = image.ID
				return
			}
		}
	}
	return
}

// ImageGetAll Image Get All.
func (d *Dao) ImageGetAll() (imageNames []string, err error) {
	var images []types.ImageSummary
	if images, err = d.dockerClient.ImageList(context.Background(), types.ImageListOptions{Filters: filters.NewArgs()}); err != nil {
		return
	}

	for _, image := range images {
		for _, tag := range image.RepoTags {
			if strings.Contains(tag, d.c.BiliHub.HostName) {
				imageNames = append(imageNames, tag)
			}
		}
	}
	return
}

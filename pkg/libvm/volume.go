package libvm

import (
	"context"
	"go-cli/pkg/libcli"
	"go-cli/pkg/libutil"
	"go-cli/pkg/libvm/vmer"
	"os"
)

func initVolumeCli(cli *libcli.GoCli) {
	// show volumes
	cli.AddCommandElem(
		nce("volume", "volume"),
		ncef("show", "show volumes", func(args []string) {
			stream, err := client.ShowVolumes(context.Background(), &vmer.VolumeMessage{})
			if err != nil {
				logger.Warn("%v", err)
				return
			}

			var streamInterface StreamInterface[*vmer.VolumeMessage] = stream

			messages, err := recvStream(streamInterface)
			if err != nil {
				logger.Warn("%v", err)
				return
			}

			libutil.PrintStructAll(messages)
		}))

	// show volume by name
	cli.AddCommandElem(
		nce("volume", "volume"),
		nce("show", "show volume by name"),
		nce("name", "volume name"),
		ncef("show", "show volume by name", func(args []string) {
			stream, err := client.ShowVolumes(context.Background(), &vmer.VolumeMessage{
				Name: args[3],
			})
			if err != nil {
				logger.Warn("%v", err)
				return
			}
			var streamInterface StreamInterface[*vmer.VolumeMessage] = stream

			messages, err := recvStream(streamInterface)
			if err != nil {
				logger.Warn("%v", err)
				return
			}

			libutil.PrintStructAll(messages)
		}))

	// create volume
	cli.AddCommandElem(
		nce("volume", "volume"),
		nce("create", "create volume"),
		nce("name", "volume name"),
		nce(libutil.NameRegex, "volume name"),
		nce("path", "volume path"),
		nce(libutil.FilePathRegex, "volume path"),
		nce("size", "volume size"),
		nce(libutil.UnitRegex, "volume size"),
		nce("origin", "volume origin"),
		ncef(libutil.NameRegex, "volume name", func(args []string) {
			_, err := client.CreateVolume(context.Background(), &vmer.VolumeMessage{
				Name: args[3],
				Path: args[5],
				Size: args[7],
				Origin: &vmer.BaseImageMessage{
					Name: args[9],
				},
			})
			if err != nil {
				logger.Warn("%v", err)
				return
			}
		}))

	// delete volume by name
	cli.AddCommandElem(
		nce("volume", "volume"),
		nce("delete", "delete volume by name"),
		nce("name", "volume name"),
		ncef(libutil.NameRegex, "volume name", func(args []string) {
			_, err := client.DeleteVolume(context.Background(), &vmer.VolumeMessage{
				Name: args[3],
			})
			if err != nil {
				logger.Warn("%v", err)
				return
			}
		}))
}

// show volumes
func (s *server) ShowVolumes(in *vmer.VolumeMessage, stream vmer.Vmer_ShowVolumesServer) error {
	var volumes []Volume

	var err error
	if in.Name == "" {
		volumes, err = vmerDB.GetAllVolumes()
	} else {
		volume, err := vmerDB.GetVolumeByName(in.Name)
		if err != nil {
			return err
		}
		volumes = append(volumes, *volume)
	}

	if err != nil {
		logger.Warn("failed to get base images: %v", err)
		return err
	}

	for _, volume := range volumes {
		if err := stream.Send(&vmer.VolumeMessage{
			Name:   volume.Name,
			Path:   volume.Path,
			Format: volume.Format,
			Size:   libutil.ConvertBytesToSize(volume.Size),
			Origin: &vmer.BaseImageMessage{
				Name: volume.Origin.Name,
			},
		}); err != nil {
			return err
		}
	}

	return nil
}

// create volume
func (s *server) CreateVolume(ctx context.Context, in *vmer.VolumeMessage) (*vmer.VolumeMessage, error) {
	// get base image by oringinal name
	baseImage, err := vmerDB.GetBaseImageByName(in.Origin.Name)
	if err != nil {
		logger.Warn("failed to get base image: %v", err)
		return in, err
	}

	// open base image file
	baseImageFile, err := os.Open(baseImage.Path)
	if err != nil {
		logger.Warn("failed to open base image file: %v", err)
		return in, err
	}

	// copy base image file to volume path
	_, err = libutil.CopyFile(baseImageFile, in.Path)
	if err != nil {
		logger.Warn("failed to copy base image file to volume path: %v", err)
		return in, err
	}

	size, err := libutil.ConvertSizeToBytes(in.Size)
	if err != nil {
		logger.Warn("failed to convert size to bytes: %v", err)
		return in, err
	}

	// resize qcow2 image
	out, err := libutil.ResizeQcow2Image(in.Path, size)
	if err != nil {
		logger.Warn("failed to resize volume file: %v, out:%s", err, out)
		return in, err
	}

	if err := vmerDB.InsertVolume(&Volume{
		Name:   in.Name,
		Path:   in.Path,
		Size:   size,
		Format: in.Format,
		Origin: *baseImage,
	}); err != nil {
		logger.Warn("failed to create volume: %v", err)
		return in, err
	}

	return in, nil
}

// delete volume
func (s *server) DeleteVolume(ctx context.Context, in *vmer.VolumeMessage) (*vmer.VolumeMessage, error) {
	// get volume by name
	volume, err := vmerDB.GetVolumeByName(in.Name)
	if err != nil {
		logger.Warn("failed to get volume: %v", err)
		return in, err
	}

	if err := vmerDB.DeleteVolume(volume); err != nil {
		logger.Warn("failed to delete volume: %v", err)
		return in, err
	}

	return in, nil
}

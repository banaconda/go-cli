package libvm

import (
	"context"
	"fmt"
	"go-cli/pkg/libcli"
	"go-cli/pkg/libutil"
	"go-cli/pkg/libvm/vmer"
	"os"
	"time"
)

func initBaseImageCli(cli *libcli.GoCli) {

	// show base image
	cli.AddCommandElem(
		nce("base-image", "base image"),
		ncef("show", "show base image", func(args []string) {
			ctx, cancel := context.WithTimeout(context.Background(), time.Second)
			defer cancel()
			stream, err := client.ShowBaseImages(ctx, &vmer.BaseImageMessage{})
			if err != nil {
				logger.Warn("%v", err)
				return
			}

			// stream to StreamInterface
			var streamInterface StreamInterface[*vmer.BaseImageMessage] = stream

			messages, err := recvStream(streamInterface)
			if err != nil {
				logger.Warn("%v", err)
				return
			}

			libutil.PrintStructAll(messages)
		}))

	// show base image by name
	cli.AddCommandElem(
		nce("base-image", "base image"),
		nce("show", "show base image"),
		nce("name", "base image name"),
		ncef(libutil.NameRegex, "", func(args []string) {
			ctx, cancel := context.WithTimeout(context.Background(), time.Second)
			defer cancel()
			stream, err := client.ShowBaseImages(ctx, &vmer.BaseImageMessage{
				Name: args[3],
			})
			if err != nil {
				logger.Warn("%v", err)
				return
			}
			var streamInterface StreamInterface[*vmer.BaseImageMessage] = stream

			messages, err := recvStream(streamInterface)
			if err != nil {
				logger.Warn("%v", err)
				return
			}

			libutil.PrintStructAll(messages)
		}))

	// upload base image
	cli.AddCommandElem(
		nce("base-image", "base image"),
		nce("upload", "upload base image"),
		nce("name", "base image name"),
		nce(libutil.NameRegex, ""),
		nce("path", "base image path"),
		ncef(libutil.FilePathRegex, "", func(args []string) {
			ctx, cancel := context.WithTimeout(context.Background(), time.Second)
			defer cancel()

			r, err := client.UploadBaseImage(ctx, &vmer.BaseImageMessage{
				Name: args[3],
				Path: args[5],
			})

			if err != nil {
				logger.Warn("%v", err)
				return
			}

			libutil.PrintStructAll(r)
		}))

	// delete base image
	cli.AddCommandElem(
		nce("base-image", "base image"),
		nce("delete", "delete base image"),
		nce("name", "base image name"),
		ncef(libutil.NameRegex, "", func(args []string) {
			ctx, cancel := context.WithTimeout(context.Background(), time.Second)
			defer cancel()

			r, err := client.DeleteBaseImage(ctx, &vmer.BaseImageMessage{
				Name: args[3],
			})

			if err != nil {
				logger.Warn("%v", err)
				return
			}

			libutil.PrintStructAll(r)
		}))

}

// show base images
func (s *server) ShowBaseImages(in *vmer.BaseImageMessage, stream vmer.Vmer_ShowBaseImagesServer) error {
	var baseImages []BaseImage

	var err error
	if in.Name == "" {
		baseImages, err = vmerDB.GetAllBaseImages()
	} else {
		baseImage, err := vmerDB.GetBaseImageByName(in.Name)
		if err != nil {
			return err
		}
		baseImages = append(baseImages, *baseImage)
	}

	if err != nil {
		logger.Warn("failed to get base images: %v", err)
		return err
	}

	for _, baseImage := range baseImages {
		if err := stream.Send(&vmer.BaseImageMessage{
			Name:   baseImage.Name,
			Path:   baseImage.Path,
			Format: baseImage.Format,
			Size:   libutil.ConvertBytesToSize(baseImage.Size),
		}); err != nil {
			return err
		}
	}

	return nil
}

// upload base image
func (s *server) UploadBaseImage(ctx context.Context, in *vmer.BaseImageMessage) (*vmer.BaseImageMessage, error) {
	// check file is exist
	if !libutil.IsExist(in.Path) {
		return nil, fmt.Errorf("file %s is not exist", in.Path)
	}

	// path is qcow2
	file, err := os.Open(in.Path)
	if err != nil {
		logger.Warn("failed to open file: %v", err)
		return nil, err
	}

	if !libutil.IsQcow2(file) {
		logger.Warn("path is not qcow2")
		return nil, fmt.Errorf("path is not qcow2")
	}

	size, err := libutil.GetFileSize(file)
	if err != nil {
		logger.Warn("failed to get file size: %v", err)
		return nil, err
	}

	// insert base image by base image message
	err = vmerDB.InsertBaseImage(&BaseImage{
		Name:   in.Name,
		Path:   in.Path,
		Format: "qcow2",
		Size:   size,
	})
	if err != nil {
		logger.Warn("failed to insert base image: %v", err)
		return in, err
	}

	return in, nil
}

// delete base image
func (s *server) DeleteBaseImage(ctx context.Context, in *vmer.BaseImageMessage) (*vmer.BaseImageMessage, error) {
	// get base image by name
	baseImage, err := vmerDB.GetBaseImageByName(in.Name)
	if err != nil {
		logger.Warn("failed to get base image: %v", err)
		return in, err
	}

	// delete base image by base image message
	err = vmerDB.DeleteBaseImage(baseImage)
	if err != nil {
		logger.Warn("failed to delete base image: %v", err)
		return in, err
	}

	return in, nil
}

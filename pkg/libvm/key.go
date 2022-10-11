package libvm

import (
	"context"
	"go-cli/pkg/libcli"
	"go-cli/pkg/libutil"
	"go-cli/pkg/libvm/vmer"
)

// init cli
func initKeyCli(cli *libcli.GoCli) {
	// show key
	cli.AddCommandElem(
		nce("key", "key"),
		ncef("show", "show key", func(args []string) {
			stream, err := client.ShowKeys(context.Background(), &vmer.KeyMessage{})
			if err != nil {
				logger.Warn("%v", err)
				return
			}

			var streamInterface StreamInterface[*vmer.KeyMessage] = stream

			messages, err := recvStream(streamInterface)
			if err != nil {
				logger.Warn("%v", err)
				return
			}

			libutil.PrintStructAll(messages)
		}))

	// show key by name
	cli.AddCommandElem(
		nce("key", "key"),
		nce("show", "show key"),
		nce("name", "show key by name"),
		ncef(libutil.NameRegex, "show key by name", func(args []string) {
			stream, err := client.ShowKeys(context.Background(), &vmer.KeyMessage{
				Name: args[3],
			})
			if err != nil {
				logger.Warn("%v", err)
				return
			}
			var streamInterface StreamInterface[*vmer.KeyMessage] = stream

			messages, err := recvStream(streamInterface)
			if err != nil {
				logger.Warn("%v", err)
				return
			}

			libutil.PrintStructAll(messages)
		}))

	// upload key
	cli.AddCommandElem(
		nce("key", "key"),
		nce("upload", "upload key"),
		nce("name", "upload key by name"),
		nce(libutil.NameRegex, "upload key by name"),
		nce("username", "upload key by username"),
		nce(libutil.NameRegex, "upload key by name"),
		nce("path", "upload key by path"),
		ncef(libutil.FilePathRegex, "upload key by name", func(args []string) {
			_, err := client.UploadKey(context.Background(), &vmer.KeyMessage{
				Name:     args[3],
				Username: args[5],
				Path:     args[7],
			})
			if err != nil {
				logger.Warn("%v", err)
				return
			}
		}))

	// delete key by name
	cli.AddCommandElem(
		nce("key", "key"),
		nce("delete", "delete key"),
		nce("name", "delete key by name"),
		ncef(libutil.NameRegex, "delete key by name", func(args []string) {
			_, err := client.DeleteKey(context.Background(), &vmer.KeyMessage{
				Name: args[3],
			})
			if err != nil {
				logger.Warn("%v", err)
				return
			}
		}))

}

// show key
func (s *server) ShowKeys(in *vmer.KeyMessage, stream vmer.Vmer_ShowKeysServer) error {
	var keys []Key

	var err error
	if in.Name == "" {
		keys, err = vmerDB.GetAllKeys()
	} else {
		key, err := vmerDB.GetKeyByName(in.Name)
		if err != nil {
			return err
		}
		keys = append(keys, *key)
	}

	if err != nil {
		logger.Warn("failed to get base images: %v", err)
		return err
	}

	for _, key := range keys {
		if err := stream.Send(&vmer.KeyMessage{
			Name:     key.Name,
			Username: key.Username,
			Key:      key.Rsa,
		}); err != nil {
			return err
		}
	}

	return nil
}

// upload key
func (s *server) UploadKey(ctx context.Context, in *vmer.KeyMessage) (*vmer.KeyMessage, error) {
	key := Key{
		Name:     in.Name,
		Username: in.Username,
		Path:     in.Path,
	}

	// read file by path
	if key.Path != "" {
		data, err := libutil.ReadFile(in.Path)
		if err != nil {
			return in, err
		}

		key.Rsa = data
	}

	if err := vmerDB.InsertKey(&key); err != nil {
		logger.Warn("failed to create key: %v", err)
		return nil, err
	}

	return in, nil
}

// delete key
func (s *server) DeleteKey(ctx context.Context, in *vmer.KeyMessage) (*vmer.KeyMessage, error) {
	// get key
	key, err := vmerDB.GetKeyByName(in.Name)
	if err != nil {
		logger.Warn("failed to get key: %v", err)
		return nil, err
	}

	if err := vmerDB.DeleteKey(key); err != nil {
		logger.Warn("failed to delete key: %v", err)
		return nil, err
	}

	return in, nil
}

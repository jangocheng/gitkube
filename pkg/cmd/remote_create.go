package cmd

import (
	"fmt"
	"io/ioutil"

	"github.com/hasura/gitkube/pkg/apis/gitkube.sh/v1alpha1"
	"github.com/hasura/gitkube/pkg/client/clientset/versioned/scheme"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

func newRemoteCreateCmd(c *Context) *cobra.Command {
	var opts remoteCreateOptions
	opts.Context = c
	var remote v1alpha1.Remote
	opts.Remote = &remote
	remoteCreateCmd := &cobra.Command{
		Use:   "create",
		Short: "Create gitkube remote on a cluster",
		RunE: func(cmd *cobra.Command, args []string) error {
			err := opts.initialize()
			if err != nil {
				return errors.Wrap(err, "creating remote failed")
			}
			err = opts.parse()
			if err != nil {
				return errors.Wrap(err, "creating remote failed")
			}
			return nil
		},
	}

	f := remoteCreateCmd.Flags()
	f.StringVarP(&opts.SpecFile, "file", "f", "", "spec file to read")

	return remoteCreateCmd
}

type remoteCreateOptions struct {
	SpecFile string
	RawData  []byte

	Remote *v1alpha1.Remote

	Context *Context
}

func (o *remoteCreateOptions) initialize() error {
	data, err := ioutil.ReadFile(o.SpecFile)
	if err != nil {
		return errors.Wrap(err, "error reading file")
	}
	o.RawData = data
	return nil
}

func (o *remoteCreateOptions) parse() error {
	gclient := o.Context.GitkubeClientSet

	d := scheme.Codecs.UniversalDeserializer()
	obj, _, err := d.Decode(o.RawData, nil, nil)
	if err != nil {
		return errors.Wrap(err, "parsing yaml as a valid remote failed")
	}
	o.Remote = obj.(*v1alpha1.Remote)

	_, err = gclient.Gitkube().Remotes(o.Remote.GetNamespace()).Create(o.Remote)
	if err != nil {
		return errors.Wrap(err, "creating remote failed")
	}

	fmt.Println(o.Remote)
	return nil
}

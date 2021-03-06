package main

import (
	"os"

	"fmt"
	"github.com/golang/glog"
	"github.com/spf13/cobra"
	"k8s.io/kops/upup/pkg/api"
	"strconv"
	"strings"
)

type GetInstanceGroupsCmd struct {
}

var getInstanceGroupsCmd GetInstanceGroupsCmd

func init() {
	cmd := &cobra.Command{
		Use:     "instancegroups",
		Aliases: []string{"instancegroup", "ig"},
		Short:   "get instancegroups",
		Long:    `List or get InstanceGroups.`,
		Run: func(cmd *cobra.Command, args []string) {
			err := getInstanceGroupsCmd.Run(args)
			if err != nil {
				glog.Exitf("%v", err)
			}
		},
	}

	getCmd.cobraCommand.AddCommand(cmd)
}

func (c *GetInstanceGroupsCmd) Run(args []string) error {
	registry, err := rootCommand.InstanceGroupRegistry()
	if err != nil {
		return err
	}

	instancegroups, err := registry.ReadAll()
	if err != nil {
		return err
	}

	if len(args) != 0 {
		m := make(map[string]*api.InstanceGroup)
		for _, ig := range instancegroups {
			m[ig.Name] = ig
		}
		instancegroups = make([]*api.InstanceGroup, 0, len(args))
		for _, arg := range args {
			ig := m[arg]
			if ig == nil {
				return fmt.Errorf("instancegroup not found %q", arg)
			}

			instancegroups = append(instancegroups, ig)
		}
	}

	if len(instancegroups) == 0 {
		fmt.Fprintf(os.Stderr, "No InstanceGroup objects found\n")
		return nil
	}

	output := getCmd.output
	if output == OutputTable {
		t := &Table{}
		t.AddColumn("NAME", func(c *api.InstanceGroup) string {
			return c.Name
		})
		t.AddColumn("ROLE", func(c *api.InstanceGroup) string {
			return string(c.Spec.Role)
		})
		t.AddColumn("MACHINETYPE", func(c *api.InstanceGroup) string {
			return c.Spec.MachineType
		})
		t.AddColumn("ZONES", func(c *api.InstanceGroup) string {
			return strings.Join(c.Spec.Zones, ",")
		})
		t.AddColumn("MIN", func(c *api.InstanceGroup) string {
			return intPointerToString(c.Spec.MinSize)
		})
		t.AddColumn("MAX", func(c *api.InstanceGroup) string {
			return intPointerToString(c.Spec.MinSize)
		})
		return t.Render(instancegroups, os.Stdout, "NAME", "ROLE", "MACHINETYPE", "MIN", "MAX", "ZONES")
	} else if output == OutputYaml {
		for _, ig := range instancegroups {
			y, err := api.ToYaml(ig)
			if err != nil {
				return fmt.Errorf("error marshaling yaml for %q: %v", ig.Name, err)
			}
			_, err = os.Stdout.Write(y)
			if err != nil {
				return fmt.Errorf("error writing to stdout: %v", err)
			}
		}
		return nil
	} else {
		return fmt.Errorf("Unknown output format: %q", output)
	}
}

func intPointerToString(v *int) string {
	if v == nil {
		return "-"
	}
	return strconv.Itoa(*v)
}

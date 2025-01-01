package herd

import (
	"encoding/json"
	"fmt"
	"sort"
	"strings"

	"golang.org/x/crypto/ssh"
)

type HostSet struct {
	hosts         []*Host
	sort          []string
	maxNameLength int
}

func NewHostSet() *HostSet {
	return &HostSet{hosts: make([]*Host, 0), sort: []string{}}
}

func (s *HostSet) Len() int {
	return len(s.hosts)
}

func (s *HostSet) Get(i int) *Host {
	return s.hosts[i]
}

func (s *HostSet) MarshalJSON() ([]byte, error) {
	return json.Marshal(s.hosts)
}

func (s *HostSet) UnmarshalJSON(b []byte) error {
	if err := json.Unmarshal(b, &s.hosts); err != nil {
		return err
	}
	s.sort = []string{}
	s.maxNameLength = maxNameLength(s.hosts)
	return nil
}

func (s *HostSet) SetSortFields(attributes []string) {
	s.sort = attributes
}

func (s *HostSet) Search(hostnameGlob string, attributes MatchAttributes) *HostSet {
	hosts := make([]*Host, 0)
	for _, host := range s.hosts {
		if host.Match(hostnameGlob, attributes) {
			hosts = append(hosts, host)
		}
	}
	return &HostSet{hosts: hosts, maxNameLength: maxNameLength(hosts)}
}

func (s *HostSet) Filter(f func(*Host) bool) *HostSet {
	hosts := make([]*Host, 0)
	for _, host := range s.hosts {
		if f(host) {
			hosts = append(hosts, host)
		}
	}
	return &HostSet{hosts: hosts, maxNameLength: maxNameLength(hosts)}
}

func (s *HostSet) AddHost(host *Host) {
	s.hosts = append(s.hosts, host)
	if l := len(host.Name); l > s.maxNameLength {
		s.maxNameLength = l
	}
}

func (s *HostSet) AddHosts(hosts *HostSet) {
	s.hosts = append(s.hosts, hosts.hosts...)
	s.maxNameLength = max(s.maxNameLength, hosts.maxNameLength)
	s.Sort()
	s.Uniq()
}

func (s *HostSet) Remove(glob string, attrs MatchAttributes) {
	newHosts := make([]*Host, 0)
	ml := 0
	for _, host := range s.hosts {
		if !host.Match(glob, attrs) {
			newHosts = append(newHosts, host)
			if l := len(host.Name); l > ml {
				ml = l
			}
		}
	}
	s.hosts = newHosts
	s.maxNameLength = maxNameLength(s.hosts)
}

func (s *HostSet) addHostKeys(allKeys []map[string][]ssh.PublicKey) {
	for _, host := range s.hosts {
		for _, set := range allKeys {
			if keys, ok := set[host.Name]; ok {
				host.publicKeys = append(host.publicKeys, keys...)
			}
		}
	}
}

func (s *HostSet) Uniq() {
	if len(s.hosts) < 2 {
		return
	}
	src, dst := 1, 0
	for src < len(s.hosts) {
		if s.hosts[src].Name != s.hosts[dst].Name {
			dst += 1
			if dst != src {
				s.hosts[dst] = s.hosts[src]
			}
		}
		src += 1
	}
	s.hosts = s.hosts[:dst+1]
}

func (s *HostSet) Sample(attributes []string, count int) {
	newHosts := make([]*Host, 0)
	buckets := make(map[string][]*Host)
host:
	for _, host := range s.hosts {
		bucket := ""
		for _, attr := range attributes {
			value, ok := host.Attributes[attr]
			if !ok {
				continue host
			}
			bucket = fmt.Sprintf("%s\000%v", bucket, value)
		}
		hosts, ok := buckets[bucket]
		if !ok {
			buckets[bucket] = []*Host{host}
		} else {
			buckets[bucket] = append(hosts, host)
		}
	}
	for _, bucket := range buckets {
		newHosts = append(newHosts, bucket[:min(count, len(bucket))]...)
	}
	s.hosts = newHosts
	s.maxNameLength = maxNameLength(s.hosts)
}

func (s *HostSet) Sort() {
	// Most common and default case
	if len(s.sort) == 0 || len(s.sort) == 1 && s.sort[0] == "name" {
		sort.Slice(s.hosts, func(i, j int) bool { return s.hosts[i].Name < s.hosts[j].Name })
	} else {
		sort.Slice(s.hosts, func(i, j int) bool {
			return s.hosts[i].less(s.hosts[j], s.sort)
		})
	}
}

func (s *HostSet) String() string {
	var ret strings.Builder
	for i, h := range s.hosts {
		if i > 0 {
			ret.WriteString(", ")
		}
		ret.WriteString(h.Name)
	}
	return ret.String()
}

func MergeHostSets(sets []*HostSet) *HostSet {
	seen := make(map[string]int)
	hosts := make([]*Host, 0)

	for _, set := range sets {
		if set == nil {
			continue
		}
		for _, host := range set.hosts {
			if existing, ok := seen[host.Name]; ok {
				hosts[existing].Amend(host)
			} else {
				seen[host.Name] = len(hosts)
				hosts = append(hosts, host)
			}
		}
	}
	return &HostSet{hosts: hosts, maxNameLength: maxNameLength(hosts)}
}

func maxNameLength(hosts []*Host) int {
	ml := 0
	for _, h := range hosts {
		if l := len(h.Name); l > ml {
			ml = l
		}
	}
	return ml
}

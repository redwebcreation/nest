package config

// todo
//
//func TestNetworkOptions_UnmarshalYAML(t *testing.T) {
//	err := yaml.Unmarshal([]byte("ipv6: true"), &NetworkOptions{})
//	assert.ErrorIs(t, err, ErrMissingIpv6Pool)
//}
//
//func TestNetworkOptions_UnmarshalYAML2(t *testing.T) {
//	var options NetworkOptions
//	err := yaml.Unmarshal([]byte("ipv6: false"), &options)
//	assert.NilError(t, err)
//
//	//for k, pool := range options.Pools {
//	//	assert.Equal(t, docker.DefaultIpv4Pools[k].Base, pool.Base, "defaults not set correctly")
//	//	assert.Equal(t, docker.DefaultIpv4Pools[k].Size, pool.Size, "defaults not set correctly")
//	//}
//}

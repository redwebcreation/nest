# Container metadata

*  nest:container=yes
*  nest:service={service.Name}
*  nest:listening_on={service.ListeningOn}
*  nest:hosts={service.Hosts, separated with commas}
*  nest:image_version:{deployment.ImageVersion}

# TODO
* `nest config <key> <value?>` 
  strategy: local, github
  url: github.com/company/infra
  path: services

Remote server must accept the server's SSH key.
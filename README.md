# provider documentation

[terraform documentation](docs/index.md)

# background to this fork

Initial found this via terraform (https://registry.terraform.io/providers/mwudka/hetznerrobot/latest) and here in github
at https://github.com/mwudka/terraform-provider-hetznerrobot, but it seems unfortunately it looks like it is no
longer maintained. Then author of https://github.com/Peters-IT/terraform-provider-hetzner-robot made a fork of what I consider to be good improvements and new features
from https://github.com/SLoeuillet/terraform-provider-hetznerrobot, and drop its maintenance.
So I made a fresh new fork to keep it up to date.

This software comes without any guarantee of functionality.

Feel free to submit merge/pull requests.

# build
## local
```
goreleaser release --snapshot --skip-sign --clean
```

## github
works with github action and goreleaser/action automatically at each new tag
package sops

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

const configTestDataSourceSopsExternal_basic = `
data "sops_external" "test_basic" {
  source     = <<EOF
hello: ENC[AES256_GCM,data:gzR9Gz4=,iv:cbMZU1nUyo5mFCW+Vel2UYbnbMA/0wKxsQzy/WVAYw8=,tag:tETDJMCJYo+4K4LwsSw4Dw==,type:str]
integer: ENC[AES256_GCM,data:9w==,iv:8gMmdZTOgdGdHlXvipvz4qchxFWMwKg95Zzvh/I84G4=,tag:jxzxI0m7stEK+zr+yc0Wsg==,type:int]
float: ENC[AES256_GCM,data:UtBf,iv:Z7hdgplz8QP9JyC/DX5WWiazdYjBvZTVolvyf9VNvyw=,tag:v8VujtMTcuM4GuZKGreVbA==,type:float]
bool: ENC[AES256_GCM,data:xW/sRw==,iv:0vXeg5/SBUDo8dmHHpDTdxMwpoCdx+ERE7dq4UgqVsc=,tag:RweVRArPBVskGtnPBiQ0Yg==,type:bool]
sops:
    kms: []
    gcp_kms: []
    azure_kv: []
    lastmodified: '2019-04-26T18:43:59Z'
    mac: ENC[AES256_GCM,data:UdHBCIrxfP+FjXwi0++Y1MUdAZ3hAa34OfG/w911zimF2YR+Mqv7PD15Osqa9GotQ5idzJEAzvz6pRVm7J388s0g2E53zBjCfLO/dcrkmVRdjTw2WYM17ewGM61HlNB9EKPe38B/eTH6PP1pTs5vjplEM/3FDblglKw8koUDdp0=,iv:LmRycuJjAoyGaY8qazR6G5CEuyD8JYCe3OO9UTek6kE=,tag:pfpB4HNE1qVmhO1QdZvVkQ==,type:str]
    pgp:
    -   created_at: '2019-01-23T10:01:20Z'
        enc: |-
            -----BEGIN PGP MESSAGE-----

            wcBMA/FdPFBXWyBuAQgAhJwnHuIY66QdnnWx2Nh6nzhBMogJtKLT4qA7ostnfMXX
            Qo1oTd5OAKT16dDZCUl8TMZqnQUzdDaVQ7H/rUOJ38EfkZBTr120JOoJuCbrUBTt
            uLEPrpgrUF8KSnRBuwnRECfU7jEi8QEwbKL4zQJREHf5I5O1iS4ZNg4h/5O7JHC/
            Hfg/pqmwN1LtEZvJZDen9CMipsO0fHqR2N4UDYuDimwIlMi0ziaq6pO0T+PlNQdn
            a9nzZwKk0pQOl80YRZQHZbSPpegOXwyFSMXKn/xfGo4YVjmpGJ5aO9ZJMLWylMgb
            VI8EzHu2ftskyGuykoOboSYAoIRb/LgGCMS4amfvmdLgAeSj/tqpz3VZdubUot9i
            fysa4epZ4FzgN+H+6OAp4lMnxKXgKOWPh3uUgyOo6wdgrJEYsZ2AQHNDjojedNzZ
            HGwPrD0VM+BW5B8PY9y7GVaNNhM+V+ZouzTio5TZpOHfLQA=
            =qDZi
            -----END PGP MESSAGE-----
        fp: 3CE5CC7219D6597CE6488BF1BF36CD3D0749A11A
    unencrypted_suffix: _unencrypted
    version: 3.2.0
EOF

  input_type = "yaml"
}`

func TestDataSourceSopsExternal(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		ProtoV5ProviderFactories: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: configTestDataSourceSopsExternal_basic,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.sops_external.test_basic", "data.hello", "world"),
					resource.TestCheckResourceAttr("data.sops_external.test_basic", "data.integer", "0"),
					resource.TestCheckResourceAttr("data.sops_external.test_basic", "data.float", "0.2"),
					resource.TestCheckResourceAttr("data.sops_external.test_basic", "data.bool", "true"),
				),
			},
		},
	})
}

const configTestDataSourceSopsExternal_nested = `
data "sops_external" "test_nested" {
  source     = <<EOF
db:
    user: ENC[AES256_GCM,data:FPeD,iv:J72gLGxfRX+8PZZrD7f5/7zPQPbMgBxL7OUxyvvFH1A=,tag:qJvyhtA4MXRjUkuYMctPlg==,type:str]
    password: ENC[AES256_GCM,data:XrL/,iv:Z/PsuhQGQVEg2ri6odjnf3aWr3U01JAz8R7MJX67Gz8=,tag:iX48abNIzKJMsxIZ3N7DKw==,type:str]
sops:
    kms: []
    gcp_kms: []
    azure_kv: []
    lastmodified: '2019-01-23T12:37:02Z'
    mac: ENC[AES256_GCM,data:SHWOM+zYaRt9e3jEiQK6bUgjEejoRm25CvWH1z1iYhvLzavmBbaCT+L0W7kw9pLZlTSXlGJfJW94sd5LFUgZ/cDvzy/IwNeU292wN9zeq0It2aTQtBWoKY+djjB3A1OoEiNoi/EBld4JfX81Jf5CCUT/LZevTawkig3URhmaIH8=,iv:WOG0Ssd0RoVpPwY8PFN7iKvSRNkmqsBIgA5xqtGs+xE=,tag:ZfV7UjxjmgXAWeWqn2u/FA==,type:str]
    pgp:
    -   created_at: '2019-01-23T10:01:20Z'
        enc: |-
            -----BEGIN PGP MESSAGE-----

            wcBMA/FdPFBXWyBuAQgAhJwnHuIY66QdnnWx2Nh6nzhBMogJtKLT4qA7ostnfMXX
            Qo1oTd5OAKT16dDZCUl8TMZqnQUzdDaVQ7H/rUOJ38EfkZBTr120JOoJuCbrUBTt
            uLEPrpgrUF8KSnRBuwnRECfU7jEi8QEwbKL4zQJREHf5I5O1iS4ZNg4h/5O7JHC/
            Hfg/pqmwN1LtEZvJZDen9CMipsO0fHqR2N4UDYuDimwIlMi0ziaq6pO0T+PlNQdn
            a9nzZwKk0pQOl80YRZQHZbSPpegOXwyFSMXKn/xfGo4YVjmpGJ5aO9ZJMLWylMgb
            VI8EzHu2ftskyGuykoOboSYAoIRb/LgGCMS4amfvmdLgAeSj/tqpz3VZdubUot9i
            fysa4epZ4FzgN+H+6OAp4lMnxKXgKOWPh3uUgyOo6wdgrJEYsZ2AQHNDjojedNzZ
            HGwPrD0VM+BW5B8PY9y7GVaNNhM+V+ZouzTio5TZpOHfLQA=
            =qDZi
            -----END PGP MESSAGE-----
        fp: 3CE5CC7219D6597CE6488BF1BF36CD3D0749A11A
    unencrypted_suffix: _unencrypted
    version: 3.2.0
EOF
  input_type = "yaml"
}`

func TestDataSourceSopsExternal_nested(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		ProtoV5ProviderFactories: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: configTestDataSourceSopsExternal_nested,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.sops_external.test_nested", "data.db.user", "foo"),
					resource.TestCheckResourceAttr("data.sops_external.test_nested", "data.db.password", "bar"),
				),
			},
		},
	})
}

const configTestDataSourceSopsExternal_raw = `
data "sops_external" "test_raw" {
  source     = <<EOF
{
	"data": "ENC[AES256_GCM,data:LOXVpbH9B6ZV+V16esP9HQ==,iv:BTeG8dZCpG1sNpnajwPdP1v/fBk/647CQ+y9ns3nn78=,tag:bPir7Bd96JgEIZ9sLnUAMQ==,type:str]",
	"sops": {
		"kms": null,
		"gcp_kms": null,
		"azure_kv": null,
		"lastmodified": "2019-01-23T13:39:15Z",
		"mac": "ENC[AES256_GCM,data:BigJKrFMd1rpPF7MmFhjmzLlTqalyMIzGtR1xvfD5ypWVWCTd/ssTQPavoZlqhZm69REQ+wf7q3ARcAgrsrIqyk49Xwnxu6VsZh32WQC8jfF5ILOFqukEb9JHXtxKvMcXpupmjnVQFSKMt2FxhtZsHmPNjzCtfqaV0sW2IGSGX8=,iv:/ubUR4GtloSsmClvbua4CIAqw6H66BGJb419A4Rbago=,tag:QnEhiKfdQNZwQpggLYUMOg==,type:str]",
		"pgp": [
			{
				"created_at": "2019-01-23T13:39:01Z",
				"enc": "-----BEGIN PGP MESSAGE-----\n\nwcBMA/FdPFBXWyBuAQgAhuNDTdTrGpzDEWxXniji6ocSLThtOI6k8ZthiuHvy9NO\ntkLE+IoxdW59XYqCoy8ejERx0jUTNwmvwO3+41c5ZXz9HOO/UCl6RuTTXSVfdcY5\nccbAWjaX0L1wyiqtRLSCzwdi8j9GDhWkiQSZ5eyjNfEHcV0IBQ/+D/YfxcWD58zw\nIjET/F+B/PsD4OjuX7m/V6jVT1/97nxfCZD8q28jzI9igloFaeBWHwslNHPIkCza\nchuot+dRfuixp/u0ndRDSZ1d731wKbi3EcnUVzsw5nQ6PJQaFgWTu0dHyrmH83TS\naVnm/nMBDPRyaRWIsCDUAsQXUf6QIok4+tTrbaZDedLgAeTzwxhtp6c8Y5A1yzNp\nG2Pd4aRN4Mng3+HPeeBU4igReMHgKuXt+awJiuxFvshpb/UoC3EL/qKMc1LfdOyc\nppFU8qp2tuBk5F39M5VKkiQXZdzak46X5LvioohDEOFTwQA=\n=IXZX\n-----END PGP MESSAGE-----",
				"fp": "3CE5CC7219D6597CE6488BF1BF36CD3D0749A11A"
			}
		],
		"unencrypted_suffix": "_unencrypted",
		"version": "3.2.0"
	}
}
EOF

  input_type = "raw"
}`

func TestDataSourceSopsExternal_raw(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		ProtoV5ProviderFactories: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: configTestDataSourceSopsExternal_raw,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.sops_external.test_raw", "raw", "Hello raw world!"),
				),
			},
		},
	})
}

const configTestDataSourceSopsExternal_simplelist = `
data "sops_external" "test_list" {
  source     = <<EOF
a_list:
- ENC[AES256_GCM,data:sxhIdQ==,iv:y2HITUKrZ/JgJT+9+UI5BDj1SMaGO0pTSvjJhsSW6r4=,tag:ZuCjI7N9WKjAshTaPYpspQ==,type:str]
- ENC[AES256_GCM,data:6yi4mg==,iv:KS4tBgQAXDfpPeWD6Ew6w1jolp4L9VY7NswmqVjivPE=,tag:9P8obopM4+BN4F9Jk2vGig==,type:str]
sops:
    kms: []
    gcp_kms: []
    azure_kv: []
    lastmodified: '2019-04-26T17:16:27Z'
    mac: ENC[AES256_GCM,data:6JzbXJfM7e4xP6qhGCSWQLal9YlXg2LmI9LSoX753Lu0uh6HuG7m02yA06Ls66jnNDzvHnxvCUobeuakuZeUNCp0+omu1pVgbs28Cx9cEya8SVwgrfBW9pQQMC8LEXSvesymDH4d78cWSUZrhLG6glOxTSZjV8Odl2/DSufuR2o=,iv:nh13+n7+ESX0eI0XmPOG9VgE9TZGz5YHjWKsNAdw3DI=,tag:Ucj6NSfEuRJIFGfaxc23Ww==,type:str]
    pgp:
    -   created_at: '2019-04-26T17:16:11Z'
        enc: |-
            -----BEGIN PGP MESSAGE-----

            wcBMA/FdPFBXWyBuAQgAmkN+YgyBOF+823IZdmGecxMWkuIB06wdRr339y1tGehi
            h1FxLlrwJU59ITCgjdzBJ0z0UCOqBP5qnwpzINu51LjtLEDMk9UOOmMfJdKLZ+5W
            3O2YuTXgKl2MfPqt6Oy17pGZaiSTNyrvI29TkkPyhi3fuTr0stg9LxL4s9qQvWjM
            kadq4ww3wwDL7VgxFxUfgF/CJALtRrdAbO3Fa63JXvOpeoa7huU75dnFleGDhFos
            WYNY2oK3U9q/wk9XtlTuArotALrveQI+UQgwQG9+19UMTTJyTcc26zXHIS7ROHND
            5qws6zlMhzKRXbvrH9CYbp1CSSvgOyG+UY2nphKURdLgAeQBS9lUL7wFklvXmgvB
            Dig84cjD4AngTuE7QeAk4h4xTlLgD+WTZFtU4gpa9061GRRGM8FGVw3DaRpXWHy0
            SUH3c/XifuCs5MExqB7LXKtX1SP81ynrTcLi0MH1r+Gt3gA=
            =O+MY
            -----END PGP MESSAGE-----
        fp: 3CE5CC7219D6597CE6488BF1BF36CD3D0749A11A
    unencrypted_suffix: _unencrypted
    version: 3.2.0
EOF
  input_type = "yaml"
}`

func TestDataSourceSopsExternal_simplelist(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		ProtoV5ProviderFactories: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: configTestDataSourceSopsExternal_simplelist,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.sops_external.test_list", "data.a_list.0", "val1"),
					resource.TestCheckResourceAttr("data.sops_external.test_list", "data.a_list.1", "val2"),
				),
			},
		},
	})
}

const configTestDataSourceSopsExternal_complexlist = `
data "sops_external" "test_list" {
  source     = <<EOF
a_list:
-   name: ENC[AES256_GCM,data:nX7D,iv:HdTElsGgx0Z2LcNmUGSbMCTyVhRD0UQMi8ztAEYQqJQ=,tag:D/sIqf+TXpUhyTbugEPpww==,type:str]
    index: ENC[AES256_GCM,data:uw==,iv:1Tz2jR86XjGVRKwaY3XO8gAH102xEKDFr/MNQbxO8GU=,tag:avkgp+Q9CraUC+hjuhYyXQ==,type:int]
-   name: ENC[AES256_GCM,data:uNO3,iv:y9ip72kPma2nTrJsqCn/+DDKLc6GeOeuEHEO2Tf2h9A=,tag:eTux2TxifhzMtWFYmlMnKA==,type:str]
    index: ENC[AES256_GCM,data:fA==,iv:6Vmcm5FRa3vLwtt5P4IjiOtJi76TfKIg+D8Bx10jaTA=,tag:EN7hJoYOXmyEynvCy7GZ6w==,type:int]
sops:
    kms: []
    gcp_kms: []
    azure_kv: []
    lastmodified: '2019-04-26T18:39:26Z'
    mac: ENC[AES256_GCM,data:Fw/zyoOVaQtGxVSdg2Wz5IHeRSuubjVb4ll8VPd6Prt4382gevlkzuv0TSAj9wAgGSiuXjeGU397kUkmDksdtsMgieh7XQxPuIoBHMaXuyvOtwBWtli8yAIkgU/lRr4Ablp3F8ZycHXPrNEm2oLonJLeSJDQKjJm4NsSP6brBs4=,iv:0JOTx4zHLoLmQFgMnh20RU1Sk0ONGR/gSoVMMHVGvFU=,tag:/BX9c/6HMM+0l38zWEwP4w==,type:str]
    pgp:
    -   created_at: '2019-04-26T18:38:53Z'
        enc: |-
            -----BEGIN PGP MESSAGE-----

            wcBMA/FdPFBXWyBuAQgAhUoPCTPjOpBexkgh5dMr2LTCb4ZsajkTXTa9a/wIJiBn
            TT1FRsQE2W+S9Yb/ClCz+ULearuUVYH0pUp7k+MDbpMt/SOMlIEA9JO0H631LqOS
            YLssnVOP/dsMH8uyhNCVuyLOHvVB3WMMxED+ic1m8oSbokqtIyCz5hmwR5MChebC
            nB42lqM8ZzRDS8DEBCykv78ityQFuLatog787sNxL9ExSeQ9iuLuu84UT4dWI4XF
            WUwwzyT3AUMbBkqftkucIi0iut+AORlgzyNAFlxxn4jXU10yl6iZvHj/Y76rJppm
            i2C4E15bS8fLrFtX7PsfnMLJOOSS+sulwr4THCFt39LgAeQtfxdz1iufUxQ+ePj4
            0w8x4ct/4EvgUuGE0OAD4oFkaCDgy+VbRWJQZn25OFLGUDMF+AUwnLnUu89hd/Ls
            ymrlCHWwIODS5M6xr6Rwx4agrWiURZUc95HiCknA7+FI3QA=
            =OUk7
            -----END PGP MESSAGE-----
        fp: 3CE5CC7219D6597CE6488BF1BF36CD3D0749A11A
    unencrypted_suffix: _unencrypted
    version: 3.2.0
EOF
  input_type = "yaml"
}`

func TestDataSourceSopsExternal_complexlist(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		ProtoV5ProviderFactories: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: configTestDataSourceSopsExternal_complexlist,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.sops_external.test_list", "data.a_list.0.name", "foo"),
					resource.TestCheckResourceAttr("data.sops_external.test_list", "data.a_list.0.index", "0"),
					resource.TestCheckResourceAttr("data.sops_external.test_list", "data.a_list.1.name", "bar"),
					resource.TestCheckResourceAttr("data.sops_external.test_list", "data.a_list.1.index", "1"),
				),
			},
		},
	})
}

const configTestDataSourceSopsExternal_json = `
data "sops_external" "test_json" {
  source     = <<EOF
{
	"hello": "ENC[AES256_GCM,data:vye/uc0=,iv:CasWaUwDHpLDkGTPrIE5Z4bI2KEBCtdw94ROfL4qlbE=,tag:I7OBTDV8JsrjLg4SfRNf8w==,type:str]",
	"integer": "ENC[AES256_GCM,data:1Q==,iv:xF2EsP5hxUpkUcS8OjsFWgSQ2D1dXxf63pnajpkFIuE=,tag:LJRnzsVNl7Ymh8yrdwdpnA==,type:float]",
	"float": "ENC[AES256_GCM,data:TK3k,iv:O64HQZG4XDATN4c3k8VTaATuSR59ynhWEQsQRUSwSog=,tag:X1j4tMpV/vlQ5xCjzuXymw==,type:float]",
	"bool": "ENC[AES256_GCM,data:QXxEsw==,iv:GVg4UD+/1VhA4QqSF6RP3YPJsdxT/1xJcg5NLzJkTzA=,tag:CLMVT+T3x8gkDsaOyUmbVw==,type:bool]",
	"sops": {
		"kms": null,
		"gcp_kms": null,
		"azure_kv": null,
		"lastmodified": "2019-08-07T23:18:36Z",
		"mac": "ENC[AES256_GCM,data:rFj23lyLaFONFQod8wlxTCCGysuCzNTQglRxqKXa5CZ9Q89jC8fnIdDUttf2oFHrfJHdvveeDbJYOoO2yEYfWr6Ty1MWzrJPyacUAGRF05PFpr0u+4xkjZNGLi5Cdg6VHb7uUu4+9EKd9d2A1bB6dWt1bEE0w3J4Il0uxn0JOMw=,iv:ERK3tzfvJzIWMcn97zg/vZ7EkWlByUUA8KczMJFPgZ0=,tag:mS+ICHZmoStgCnyIFgvbWQ==,type:str]",
		"pgp": [
			{
				"created_at": "2019-08-07T23:18:36Z",
				"enc": "-----BEGIN PGP MESSAGE-----\n\nhQEMA/FdPFBXWyBuAQf/eGZpJCY8TJ720XSH5rscUv19MC+C7v+xugWXdaBpNkt1\nsdZ4iJUpWFRv+ofYMm607AbhCfHRYTZtP7EiVGVl4yKwd+ztWgHSwXwQ/4WKx6QT\nWckxQlzRjxMhiJJuWsRmnz92PcZsb8yY7AsupPi9RaCykTVe4Fnx6xAdtA4l92n+\n6DQbVzFmfH/LOXJZK1YeFTeZoKiK+SJ+gMqcwoefy17F+fyfu6PWxWuUknDamReh\nNkPV6cbOtNl/J9+khWVlZObZ/DUCilOGkcF5H3qjOOdPgHyjRVRe3JUwb1uL6gBW\nsJRS54kQoHk3C68ZgxMc9GrXCsFv1jdZKlMkvb+Xo9JeAZVGiHNs4TnkGAuMAt49\nYnLx9tvpwTrMmZCbZRYV3jwl4TZDqklQ1kY5Qd96fun0KYiz1l2MeRxs3EhS2gLF\nABv8YNlQ7a7uwpWcqZtx1/Cwmdrc/EjK5fmydww9aA==\n=ShdQ\n-----END PGP MESSAGE-----\n",
				"fp": "3CE5CC7219D6597CE6488BF1BF36CD3D0749A11A"
			}
		],
		"unencrypted_suffix": "_unencrypted",
		"version": "3.3.1"
	}
}
EOF
  input_type = "json"
}`

func TestDataSourceSopsExternal_json(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		ProtoV5ProviderFactories: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: configTestDataSourceSopsExternal_json,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.sops_external.test_json", "data.hello", "world"),
					resource.TestCheckResourceAttr("data.sops_external.test_json", "data.integer", "0"),
					resource.TestCheckResourceAttr("data.sops_external.test_json", "data.float", "0.2"),
					resource.TestCheckResourceAttr("data.sops_external.test_json", "data.bool", "true"),
				),
			},
		},
	})
}

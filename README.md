# mqttauth

A simple Go program that acts as an authenticator for HMQ.

I use [hmq](https://github.com/fhmq/hmq) as my MQTT broker running inside [gokrazy](https://gokrazy.org/). To enable authentication across my MQTT network, I created this program to act as a gateway authenticator. It allows hmq to verify that MQTT clients are properly authorized to publish and subscribe.

At the time of writing, hmq also needs to be modified, as it does not provide an easy way to configure the authenticator endpoint. To change it, navigate to `hmq/plugins/auth/authhttp` and open the `http.json` file. Inside this file, apply the following configuration:

``` json
{
    "auth": "http://<ip of the authenticator>:9393/mqtt/auth",
    "acl": "http://<ip of the authenticator>:9393/mqtt/acl",
    "super": "http://127.0.0.1:9090/mqtt/superuser"
}
```

My `hmq.json` configuration file looks like this:

``` json
{
	"port": "8090",
	"host": "0.0.0.0",
	"plugins": {
		"auth": "authhttp"
	}
}
```

And my `config.json` in gokrazy:

``` json
"github.com/fhmq/hmq": {
	"ExtraFilePaths":
	{
	"/etc/hmq/hmq.json": "hmq.json",
	"/etc/hmq/http.json":"http.json"
	},
	"CommandLineFlags":
	[
		"--c=/etc/hmq/hmq.json"
	],
	WaitForClock": true
},

"github.com/BrunoTeixeira1996/mqttauth":{
	"CommandLineFlags":
	[
		"-mqtt_username=<username>",
		"-mqtt_password=<password>"
	]
}
```


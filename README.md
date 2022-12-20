# Kotsd

A cli to run commands against **multiple** kots instances.

Note: This is not officially supported by Replicated. Review the usage instructions and Security concerns before using.

## Usage

### Config

The following commands will only update the local config. Use these to add kots instances to the yaml config. Or list all instances in the config. None of these commands will connect with a running kotsadm instance.

* `kotsd add-instance`: Add a kots instance to the config file.
  * Flags:
  ```
  -e, --endpoint string   URL of the kots instance, for example http://10.10.10.5:8800
  -n, --name string       Name of the kots instance (should be unique)
  -v, --tlsVerify         If false, insecure or self signed tls for the kots instance will be allowed (default true)
  ```
  * Examples:
  ```
  kotsd add-instance --name kurl --endpoint https://35.231.189.178:8800 --tlsVerify=0 --config ./.josh/kots.yaml
  kotsd add-instance --name gke --endpoint http://localhost:8800 --config ./.josh/kots.yaml
  ```
* `kotsd list-instances`: List all kots instances from the config file.
  * Examples:
  ```
  kotsd list-instances --config ./.josh/kots.yaml
  ```

### Runtime

The following commands will connect with the specified kotsadm instance(s) using the configuration from the config file, and execute the corresponding command. If no named instances are specified, the command will be executed for all instances in the config.

* `kotsd list [...name]`: List the current kots and application version for each instance. Also list if any new version is available.
  * Examples:
  ```
  kotsd list --config ./.josh/kots.yaml
  kotsd list gke kurl --config ./.josh/kots.yaml
  ```
* `kotsd update [...name]`: Update all the apps on the instance to the latest sequence number available.
  * Examples:
  ```
  kotsd update --config ./.josh/kots.yaml
  kotsd update gke --config ./.josh/kots.yaml
  ```

### Backlog


* kotsd delete-instance --name=
* kotsd update-instance --name= [--endpoint=] [--password=] --tlsVerify
* kotsd version: Display the version of the kotsd binary

* kotsd update [name1.0, name2.0, name3.0, ...]: Update only the first (#0) app on the instance to the latest version available.
* kotsd redeploy [name1, name2, name3, ...]
* kotsd remove [name1, name2, name3, ...]
* kotsd set-config [name1, name2, name3, ...]
* kotsd support-bundle [name1, name2, name3, ...]

## Security

* When the `kots` instance configuration is saved, it will be base64 encoded (not secure)!

Conclusion: When using this cli, you have to be 100% sure that the configuration is **only** accessible by you!

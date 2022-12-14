# Kotsd

A cli to run commands against **multiple** kots instances.

Note: This is not officially supported by Replicated. Review the usage instructions and Security concerns before using.

![Test workflow](https://github.com/jdewinne/kotsd/actions/workflows/test.yaml/badge.svg)
![Release workflow](https://github.com/jdewinne/kotsd/actions/workflows/release.yaml/badge.svg)
![Snyk workflow](https://github.com/jdewinne/kotsd/actions/workflows/snyk.yaml/badge.svg)
[![Maintainability](https://api.codeclimate.com/v1/badges/2bbb51a1cb5240a2ccd8/maintainability)](https://codeclimate.com/github/jdewinne/kotsd/maintainability)
[![Test Coverage](https://api.codeclimate.com/v1/badges/2bbb51a1cb5240a2ccd8/test_coverage)](https://codeclimate.com/github/jdewinne/kotsd/test_coverage)

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
* `kotsd update-instance`: Update a kots instance in the config file.
  * Flags:
  ```
  -e, --endpoint string   URL of the kots instance, for example http://10.10.10.5:8800
  -n, --name string       Name of the kots instance (should be unique)
  -v, --tlsVerify         If false, insecure or self signed tls for the kots instance will be allowed (default true)
  ```
  * Examples:
  ```
  kotsd update-instance --name kurl --endpoint https://35.231.189.178:8800 --tlsVerify=0 --config ./.josh/kots.yaml
  kotsd update-instance --name gke --endpoint http://localhost:8800 --config ./.josh/kots.yaml
  ```
* `kotsd list-instances`: List all kots instances from the config file.
  * Examples:
  ```
  kotsd list-instances --config ./.josh/kots.yaml
  ```
* `kotsd delete-instance`: Delete a kots instance from the config file.
  * Flags:
  ```
  -n, --name string       Name of the kots instance (should be unique)
  ```
  * Examples:
  ```
  kotsd delete-instance --name gke
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
  * Flags:
  ```
  -s, --slug string   Specify the application slug of the application you wish to update
  ```
  * Examples:
  ```
  kotsd update --config ./.josh/kots.yaml
  kotsd update gke --config ./.josh/kots.yaml
  kotsd update gke kurl -s sentry-pro --config ./.josh/kots.yaml
  ```
* `kotsd redeploy [...name]`: Redeploy all the apps on the instance.
  * Flags:
  ```
  -s, --slug string   Specify the application slug of the application you wish to redeploy
  ```
  * Examples:
  ```
  kotsd redeploy --config ./.josh/kots.yaml
  kotsd redeploy gke --config ./.josh/kots.yaml
  kotsd redeploy gke kurl -s sentry-pro --config ./.josh/kots.yaml
  ```
* `kotsd remove [...name]`: Remove all the apps on the instance.
  * Flags:
  ```
  -s, --slug string   Specify the application slug of the application you wish to remove
  ```
  * Examples:
  ```
  kotsd remove --config ./.josh/kots.yaml
  kotsd remove gke --config ./.josh/kots.yaml
  kotsd remove gke kurl -s sentry-pro --config ./.josh/kots.yaml
  ```

### Backlog

* kotsd set-config [name1, name2, name3, ...]
* kotsd support-bundle [name1, name2, name3, ...]
* kotsd check-for-update [name1, name2, name3, ...]
* kotsd set-update-cron-schedule --pattern * * * * * [name1, name2, name3, ...]
* kotsd set-auto-update [on, off, patch, minor, ...] [name1, name2, name3, ...]

## Security

* When the `kots` instance configuration is saved, it will be base64 encoded (not secure)!

Conclusion: When using this cli, you have to be 100% sure that the configuration is **only** accessible by you!

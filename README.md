# Kotsd

A cli to run commands against **multiple** kots instances.

Note: This is not officially supported by Replicated. Review the usage instructions and Security concerns before using.

## Usage

### Config

The following commands will only update the local config. Use these to add kots instances to the yaml config. Or list all instances in the config. None of these commands will connect with a running kotsadm instance.

* kotsd add-instance --name= --endpoint= --tlsVerify
* kotsd list-instances

### Runtime

The following commands will connect with the specified kotsadm instance(s) using the configuration from the config file, and execute the corresponding command. If no named instances are specified, the command will be executed for all instances in the config.

* `kotsd list [name1, name2, name3, ...]`: List the current kots and application version for each instance. Also list if any new version is available.
* `kotsd update [name1, name2, name3, ...]`: Update the instance to the latest version available.

### Backlog

kotsd delete-instance --name=
kotsd update-instance --name= [--endpoint=] [--password=] --tlsVerify

kotsd redeploy [name1, name2, name3, ...]
kotsd remove [name1, name2, name3, ...]
kotsd set-config [name1, name2, name3, ...]
kotsd support-bundle [name1, name2, name3, ...]

## Security

* When the `kots` instance configuration is saved, it will be base64 encoded (not secure)!

Conclusion: When using this cli, you have to be 100% sure that the configuration is **only** accessible by you!

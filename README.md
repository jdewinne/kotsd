struct instances
  name
  endpoint
  base64 encoded password

Stored [] instances

kotsd add-instance --name= --endpoint= --password= --tlsVerify
kotsd delete-instance --name=
kotsd update-instance --name= [--endpoint=] [--password=] --tlsVerify
kotsd list-instances
  output:
  name | appversion | connected
kotsd update [name1, name2, name3, ...]


Backlog
kotsd redeploy [name1, name2, name3, ...]
kotsd remove [name1, name2, name3, ...]
kotsd set-config [name1, name2, name3, ...]
kotsd support-bundle [name1, name2, name3, ...]
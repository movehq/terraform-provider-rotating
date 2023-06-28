# Blue Green Rotating Demo

This demo shows how to use the rotating provider to do a blue green deployment. The provider gives you an active
set and you can rotate something (like access keys) in an inactive set.

In this example, we will use a `null_resource` but I will talk about a real world example (rotating access keys).

Say we have a deployment in kubernetes that needs to access AWS resources and we want to rotate those keys
every thirty days. We can use the rotating provider to do this.

## Step 1: Create the active and inactive sets

Using the provided example, we will initially create both the active an inactive sets. 
The active set will be the set that is currently used in our fictitious deployment. 
The inactive set will be the set that we will rotate to, later.

After first applying the provided configuration, you should see the following output:

```bash
╰─ terraform apply

Terraform used the selected providers to generate the following execution plan. Resource actions are indicated with the following symbols:
  + create

Terraform will perform the following actions:

  # null_resource.blue will be created
  + resource "null_resource" "blue" {
      + id = (known after apply)
    }

  # null_resource.green will be created
  + resource "null_resource" "green" {
      + id = (known after apply)
    }

  # random_pet.blue will be created
  + resource "random_pet" "blue" {
      + id        = (known after apply)
      + keepers   = {
          + "uuid" = (known after apply)
        }
      + length    = 2
      + separator = "-"
    }

  # random_pet.green will be created
  + resource "random_pet" "green" {
      + id        = (known after apply)
      + keepers   = {
          + "uuid" = (known after apply)
        }
      + length    = 2
      + separator = "-"
    }

  # rotating_blue_green.this will be created
  + resource "rotating_blue_green" "this" {
      + active               = (known after apply)
      + blue_uuid            = (known after apply)
      + green_uuid           = (known after apply)
      + id                   = (known after apply)
      + rotate_after_minutes = 2
      + rotate_timestamp     = (known after apply)
    }

Plan: 5 to add, 0 to change, 0 to destroy.

Changes to Outputs:
  + active = (known after apply)

Do you want to perform these actions?
  Terraform will perform the actions described above.
  Only 'yes' will be accepted to approve.

  Enter a value: yes

rotating_blue_green.this: Creating...
rotating_blue_green.this: Creation complete after 0s [id=7de78b1c-90ee-489c-ac61-10083f26e725]
random_pet.green: Creating...
random_pet.blue: Creating...
random_pet.blue: Creation complete after 0s [id=deep-pup]
random_pet.green: Creation complete after 0s [id=sure-monkey]
null_resource.blue: Creating...
null_resource.blue: Creation complete after 0s [id=5572471609254710276]
null_resource.green: Creating...
null_resource.green: Creation complete after 0s [id=667576671584373708]

Apply complete! Resources: 5 added, 0 changed, 0 destroyed.

Outputs:

active = "blue"
```


We can now save the active "blue" deployment to a secret and have our kubernetes deployment use it.

## Step 2: Rotate the inactive set

After the predetermined time has passed, we can rotate the inactive set's credentials and make it the active set.

You will notice that we rotate the inactive set and then make it active. This is so we can give our kubernetes deployment
time to pick up the new secret that contains the new credentials and we dont rip out the rug from under it.

```bash
╰─ terraform apply
rotating_blue_green.this: Refreshing state... [id=7de78b1c-90ee-489c-ac61-10083f26e725]
random_pet.green: Refreshing state... [id=sure-monkey]
random_pet.blue: Refreshing state... [id=deep-pup]
null_resource.blue: Refreshing state... [id=5572471609254710276]
null_resource.green: Refreshing state... [id=667576671584373708]

Terraform used the selected providers to generate the following execution plan. Resource actions are indicated with the following symbols:
-/+ destroy and then create replacement

Terraform will perform the following actions:

  # null_resource.green will be replaced due to changes in replace_triggered_by
-/+ resource "null_resource" "green" {
      ~ id = "667576671584373708" -> (known after apply)
    }

  # random_pet.green must be replaced
-/+ resource "random_pet" "green" {
      ~ id        = "sure-monkey" -> (known after apply)
      ~ keepers   = { # forces replacement
          ~ "uuid" = "82daadb8-a0bc-4ec6-bade-a8e1ff7e9077" -> "72ecc21d-01cb-4d4f-9e90-63ad592571ba"
        }
        # (2 unchanged attributes hidden)
    }

Plan: 2 to add, 0 to change, 2 to destroy.

Changes to Outputs:
  ~ active = "blue" -> "green"

Do you want to perform these actions?
  Terraform will perform the actions described above.
  Only 'yes' will be accepted to approve.

  Enter a value: yes

null_resource.green: Destroying... [id=667576671584373708]
null_resource.green: Destruction complete after 0s
random_pet.green: Destroying... [id=sure-monkey]
random_pet.green: Destruction complete after 0s
random_pet.green: Creating...
random_pet.green: Creation complete after 0s [id=assuring-walrus]
null_resource.green: Creating...
null_resource.green: Creation complete after 0s [id=2757211547200132632]

Apply complete! Resources: 2 added, 0 changed, 2 destroyed.

Outputs:

active = "green"
```

## Rotate Forever And Stay Up

You can now rotate your credentials forever and never have to worry about downtime.
As you can see above, we are only rotating inactive credentials.
The benefit here is that we are giving our deployments time to rotate away from the keys before we destroy them.
As long as we pull in that new secret and restart our deployments before the next rotation, we will never have downtime.


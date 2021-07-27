# go-mailchimp
This is a simple package that wraps the v3 MailChimp Marketing API. The current usage of this module is quite specific, hence the lack of support for the entire marketing API. The supported features include CRUD operations on lists (or audiences) and batch creation/updating of members within a particular list. 

## Creating a client 
The client is a data structure with receiver functions for communicating with the MailChimp Marketing API. To create a new client, all you need is the API key you wish to use as well as the region of your MailChimp account. After the client has been created, it does not require any closing as it does not keep a constant connection to MailChimp. Rather, it sends separate HTTP requests for each operation the client performs. 
```go
chimp := mailchimp.NewClient("secret-key", "region")
```
For information regarding how to generate an API key and find your region, please refer to the MailChimp documentation.

## Ping MailChimp
To make sure that the client is properly set up, you can use the `Ping` receiver function. This returns a boolean and an error. If the returned boolean has a value of `true` then you are all set to go! Otherwise, there may be a problem with internet connectivity or the API key. 

```go
chimp := mailchimp.NewClient("secret-key", "region")
if up, _ := chimp.Ping(); !up {
    return
}
```

## Creating a list (audience)
To create a list, start with initialising a list builder like so:
```go
builder := mailchimp.ListBuilder{}
```
You can then use chaining receiver functions on the builder to set the parameters of the new list. There are a few required fields for lists, if these are not specified when `Build` is called, then an error will be returned. The required fields are listed below.

* `builder.Name("The name of your list")`
* `builder.PermissionReminder("The permission reminder for your list")`

Moreover, a list must also contain a contact and campaign defaults. These are instantiated manually. One example for each is given below, together with how to wire them up to the list builder. If any of the fields marked with a *required* comment are not specified when `Build` is called, then an error will be returned. 

```go
contact := mailchimp.Contact{
    Address1: "Company address",    // required
    Zip: "123 45",                  // required
    City: "Company city",           // required
    State: "Company state",         // required
    Country: "Company country",     // required
    Address2: "Secondary address",
    Company: "Company name",        // required
}

campaignDefaults := mailchimp.CampaignDefaults{
    FromName: "Name",               // required
    FromEmail: "your@email.com",    // required
    Subject: "A subject",           // required
    Language: "Some language",      // required
}

list, err := mailchimp.ListBuilder{}.
    Name("Test").
    PermissionReminder("This is a test").
    Contact(contact).
    CampaignDefaults(campaignDefaults).
    Build()
```

After a list has been built with the list builder, it can be sent through the client to MailChimp in the `CreateList` receiver function. This function will return an error either if the MailChimp response is an error or if the MailChimp response could not be unmarshalled. If no error occurs, then the newly created list will be returned to the caller with some additional information filled in by MailChimp, such as the ID of the list. 

```go
list, err := mailchimp.ListBuilder{}.[...].Build()
if err != nil {
    return handleErr(err)
}
chimp := mailchimp.NewClient("key", "region")
createdList, err := chimp.CreateList(list)
if err != nil {
    return handleErr(err)
}
```

## Fetching lists
It is possible to fetch all the lists on your MailChimp account using the client receiver function `FetchLists`. This returns a slice of List structs or an error if MailChimp responds with an error or if the response could not be unmarshalled. 

```go
chimp := mailchimp.NewClient("key", "region")
lists, err := chimp.FetchLists()
```

## Fetching a single list
It is also possible to fetch a single list, given that the ID of is already known. This can be achieved using the `FetchList` client receiver function. This function returns both a list and an error, but the error will only be non-nil if an error was returned from the MailChimp Marketing API or if there was an issue in unmarshalling the response. 

```go
chimp := mailchimp.NewClient("key", "region")
list, err := chimp.FetchList("some-list-id")
```

## Updating an existing list
To update the information regarding an existing list, the clients `UpdateList` receiver function can be used. This, of course, requires knowledge of the lists ID. Even though you can call `UpdateList` directly with a list struct, it would be advisable to first fetch the list from MailChimp, perform the necessary modifications, and then use that list object to perform the update. The suggested flow is shown below. Note that the updated list will be returned from the `UpdateList` call together with a potential error.

```go
chimp := mailchimp.NewClient("key", "region")
list, err := chimp.FetchList("some-id")
if err != nil {
    return handleErr(err)
}
list.Name = "New and improved name"
updatedList, err := chimp.UpdateList("some-id", list)
```

## Deleting a list
In order to delete a list from your MailChimp account, you must know the ID of the list beforehand. Once the ID has been acquired, the clients `DeleteList` receiver function can be invoked to perform the action. This function returns an error if an error was returned from MailChimp. 

```go
chimp := mailchimp.NewClient("key", "region")
err := chimp.DeleteList("some-id")
```

## Adding members to a list
There are two ways in which members can be added to a list. Both are described below, but first we will cover how to create new member structs. 

### Creating a member
To create a member, there is another builder that can be used. Much like the `ListBuilder`, the `MemberBuilder` has a required field, namely `EmailAddress`. For MailChimp merge fields, such as `FNAME` and `PHONE`, one can use the `MergeField` receiver function on the `MemberBuilder`. A simple example of its usage is shown below.

```go
member, err := mailchimp.MemberBuilder{}.
    EmailAddress("test@test.com").
    StatusSubscribed().
    MergeField("FNAME", "Test").
    MergeField("LNAME", "Smith").
    Build()
```

The available statuses for members are listed as the corresponding receiver function below.

* `StatusSubscribed()`
* `StatusUnsubscribed()`
* `StatusPending()`
* `StatusCleaned()`

### `Batch`
Using `Batch` to add members will only work if all the members are new. Meaning, you cannot update an existing member if the `Batch` function is used. The prerequisite knowledge to use `Batch` is the ID of the list that the members should be added to as well as the members that should be added. Please note that a maximum of **500** members can be batched for a single request as per MailChimps' specifications, if any more than that is sent to `Batch` then an error will be returned. A simple usage example for the `Batch` function is shown below.

```go
chimp := mailchimp.NewClient("key", "region")
members := createMembers()
if err := chimp.Batch("some-id", members); err != nil {
    return handleErr(err)
}
```

### `BatchWithUpdate`
The `BatchWithUpdate` function is very similar to `Batch`, with the difference being that `BatchWithUpdate` will update already existing members of the MailChimp list. Hence, if a member was subscribed with a `Batch` call, then if the same email address is found with a `BatchWithUpdate` call but with a status of `unsubscribed` then the member will be unsubscribed from the list. 

```go
chimp := mailchimp.NewClient("key", "region")
members := createMembers()
if err := chimp.BatchWithUpdate("some-id", members); err != nil {
    return handleErr(err)
}
```
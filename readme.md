# go-mailchimp
This is a simple package that wraps the v3 MailChimp Marketing API. The current usage of this module is quite specific, hence the lack of support for the entire marketing API. The supported features include CRUD operations on lists (or audiences) and batch creation/updating of members within a particular list. 

## Creating a client 
The client is a data structure with receiver functions for communicating with the MailChimp Marketing API. To create a new client, all you need is the API key you wish to use as well as the region of your MailChimp account. After the client has been created, it does not require any closing as it does not keep a constant connection to MailChimp. Rather, it sends separate HTTP requests for each operation the client performs. 
```go
chimp := mailchimp.NewClient("key", "region")
```
For information regarding how to generate an API key and find your region, please refer to the MailChimp documentation.

## Ping MailChimp
To make sure that the client is properly set up, you can use the `Ping` receiver function. This returns an error if something went wrong whilst sending the **ping**, examples of what could go wrong is a loss in internet connectivity or an invalid API key. If the returned error is `nil`, then everything is ready to go with the client.

```go
chimp := mailchimp.NewClient("key", "region")
if err := chimp.Ping(); err != nil {
    return handleErr(err)
}
```

## Creating a list (audience)
To create a list, start with initialising a list builder like so:
```go
builder := mailchimp.ListBuilder{}
```
You can then use chaining receiver functions on the builder to set the parameters of the new list. There are a few required fields for lists, if these are not specified when `Build` is called, then an error will be returned. It is for this reason that using the `ListBuilder` is advised, it will ensure that all the required data has been filled in before contacting MailChimp. The required fields are listed below.

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
It is also possible to fetch a single list, given that its ID is known beforehand. This can be achieved using the `FetchList` client receiver function. This function returns both a list and an error, but the error will only be non-nil if an error was returned from the MailChimp Marketing API or if there was an issue in unmarshalling the response. 

```go
chimp := mailchimp.NewClient("key", "region")
list, err := chimp.FetchList("list-id")
```

## Updating an existing list
To update the information regarding an existing list, the clients `UpdateList` receiver function can be used. This, of course, requires knowledge of the lists ID. Even though you can call `UpdateList` directly with a list struct, it would be advisable to first fetch the list from MailChimp, perform the necessary modifications, and then use that list object to perform the update. The suggested flow is shown below. Note that the updated list will be returned from the `UpdateList` call together with a potential error.

```go
chimp := mailchimp.NewClient("key", "region")
list, err := chimp.FetchList("list-id")
if err != nil {
    return handleErr(err)
}
list.Name = "New and improved name"
updatedList, err := chimp.UpdateList("list-id", list)
```

## Deleting a list
In order to delete a list from your MailChimp account, you must know the ID of the list beforehand. Once the ID has been acquired, the clients `DeleteList` receiver function can be invoked to perform the action. This function returns an error if an error was returned from MailChimp. 

```go
chimp := mailchimp.NewClient("key", "region")
err := chimp.DeleteList("list-id")
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

* `builder.StatusSubscribed()`
* `builder.StatusUnsubscribed()`
* `builder.StatusPending()`
* `builder.StatusCleaned()`

### `Batch`
Using `Batch` to add members will only work if all the members are new. Meaning, you cannot update an existing member if the `Batch` function is used. The prerequisite knowledge to use `Batch` is the ID of the list that the members should be added to as well as the members that should be added. Please note that a maximum of **500** members can be batched for a single request as per MailChimps' specifications, if any more than that is sent to `Batch` then an error will be returned. A simple usage example for the `Batch` function is shown below.

```go
chimp := mailchimp.NewClient("key", "region")
members := createMembers()
if err := chimp.Batch("list-id", members); err != nil {
    return handleErr(err)
}
```

### `BatchWithUpdate`
The `BatchWithUpdate` function is very similar to `Batch`, with the difference being that `BatchWithUpdate` will update already existing members of the MailChimp list. Hence, if a member was subscribed with a `Batch` call, then if the same email address is found with a `BatchWithUpdate` call but with a status of `unsubscribed` then the member will be unsubscribed from the list. 

```go
chimp := mailchimp.NewClient("key", "region")
members := createMembers()
if err := chimp.BatchWithUpdate("list-id", members); err != nil {
    return handleErr(err)
}
```

## Archiving a member from a list
First of all, make sure that you actually want to delete the member and not unsubscribe them.

To archive a member, you need to know the list ID and member email address beforehand. The archiving is simple and straight forward as shown below. 

```go
chimp := mailchimp.NewClient("key", "region")
if err := chimp.ArchiveMember("list-id", "test@test.com"); err != nil {
	return handleErr(err)
}
```

## Fetching a members tags 
It is possible to fetch all the tags associated with a given member for a given list. However, it is required that the lists ID and the members email address is known beforehand. To fetch the tags, simply use the `FetchMemberTags` receiver function on your `mailchimp.Client`. As example is given below. Please note that this function will only return an error is something went wrong on the MailChimp API side.

```go
chimp := mailchimp.NewClient("key", "region")
tags, err := chimp.FetchMemberTags("list-id", "member@email.com")
if err != nil {
    return handleErr(err)
}
```

## Adding/removing Tags
In order to add or remove tags for a given member of a given list, the members email address and the lists ID must be known beforehand. Before the operation can be performed, the client must first build a set of tags to either add or remove. This is simply done with `TagBuilder` as shown below.

```go
tag, err := mailchimp.TagBuilder{}.
    Name("my-tag").
    StatusActive().
    Build()
```

The `Build` receiver function will return an error if the tags name has not been properly specified. There are also two possible statuses for a tag, `active` and `inactive`. Their corresponding builder receiver functions are listed below.

* `builder.StatusActive()`
* `builder.StatusInactive()`

Setting a tags status to `inactive` means that if the tag already exists on the member for the specified list at MailChimp, then that tag will be removed from the member. Setting the status as `active` simply means that the tag will be added to the member. 

Once the tags has been created, you can send them to MailChimp like so:
```go
chimp := mailchimp.NewClient("key", "region")
tags := createTags()
if err := chimp.UpdateMemberTags("list-id", "member@email.com", tags); err != nil {
    return handleErr(err)
}
```

There is also another version of `UpdateMemberTags` called `UpdateMemberTagsSync`. Using `UpdateMemberTagsSync` will make sure that any automations at MailChimp based on tags are **not** ran during the update. Please note that this also means that using `UpdateMemberTags` to update the tags will cause these automations to run, if any are set up. Please note that both of these receiver functions will only return an error if one occured on the MailChimp API side.

## Webhooks
It is possible to add Webhooks unto your MailChimp audience using the `mailchimp.Client`. To create a new client, simply call `mailchimp.NewClient` with the API key and region for your MailChimp account. After creating a client you can do add, fetch and delete Webhooks on your MailChimp audience. Each of these operations are described with examples below. 

### Create Webhook
To add a Webhook you must know the URL to which events should be sent and the list ID of the audience you want to listen to. You can specify the types of events that should be sent to the Webhook using `mailchimp.WebhookEvents`, as well as the preferred event sources using `mailchimp.WebhookSources`. It is recommended to use `mailchimp.WebhookBuilder{}` when creating your Webhook as it will perform some basic validation on the data before it is sent off to MailChimp. A simple example of creating a Webhook is shown below. As can be seen, the `mailchimp.Webhook` only needs to be sent as a parameter to `chimp.CreateWebhook` in order to finalize the operation.

**Note:** Even though MailChimp uses HTTP POST requests to forward events to the Webhook, an initial HTTP GET request will be made to it at the point of creating the Webhook. Make sure to have a handler for a GET request that returns `200 OK` as well as a POST handler for your Webhook.

```go
wh, err := mailchimp.WebhookBuilder{}.
	URL("http://your-url.com/webhook").
	ListID("list-id").
	Events(mailchimp.WebhookEvents{
		Subscribe: true,
	}).
	Sources(mailchimp.WebhookSources{
		Admin: true,
	}).
	Build()
if err != nil {
	handleErr(err)
}

chimp := mailchimp.NewClient("list-id", "region")
wh, err = chimp.CreateWebhook(wh)
if err != nil {
	handleErr(err)
}
```

The full set of options for events and sources are shown below. 

```go
type WebhookEvents struct {
	Subscribe   bool
	Unsubscribe bool
	Profile     bool
	Cleaned     bool
	UpEmail     bool
	Campaign    bool
}

type WebhookSources struct {
	User  bool
	Admin bool
	API   bool
}
```

### Fetch a Webhook by ID 
It is rather straight forward to fetch a Webhook by its ID. Simply call `FetchWebhook` with the list ID and Webhook ID as parameters as shown below.

```go
chimp := mailchimp.NewClient("key", "region")
webhook, err := chimp.FetchWebhook("list-id", "webhook-id")
if err != nil {
    handleErr(err)
}
```

### Delete a Webhook
To delete a Webhook, all that is required is the list ID as well as the Webhook ID. Simply call `DeleteWebhook` and check for an error to complete the operation. An example is shown below. 

```go
chimp := mailchimp.NewClient("key", "region")
if err := chimp.DeleteWebhook("list-id", "webhook-id"); err != nil {
    handleErr(err)
}
```

## Testing
### Mocking the MailChimp provider
While running automated tests, it is very likely that you do not want `go-mailchimp` to send real requests to the MailChimp Marketing API. To avoid this, one can use the `mailchimp.NewCustomDependencyClient` to instantiate a client in place of the `mailchimp.NewClient` function. This function requires a value of the type `mailchimp.MailChimpProviderMock` to be sent in as a parameter. Using this mock, you can define the behaviour of the MailChimp endpoints for `GET`, `PATCH`, `POST` and `DELETE` calls. Thus, if you need to test how your software behaves when an error is returned from `go-mailchimp` you can simply define a function that returns an arbitrary error. By inspecting for example the `PostCalls` field on the `mailchimp.MailChimpProviderMock` you can also see how many `POST` requests were made during the test. 

The `mailchimp.MailChimpProviderMock` struct is specified below.

```go
type MailChimpProviderMock struct {
	PostMock    func(uri string, payload interface{}) ([]byte, error)
	PostCalls   int
	GetMock     func(uri string) ([]byte, error)
	GetCalls    int
	PatchMock   func(uri string, payload interface{}) ([]byte, error)
	PatchCalls  int
	DeleteMock  func(uri string) ([]byte, error)
	DeleteCalls int
}
```

An example of using the mock is shown below.

```go
mock := mailchimp.MailChimpProviderMock{
	PostMock: func(s string, i interface{}) ([]byte, error) {
		return nil, errors.New("something went wrong")
	},
}
chimpMock := mailchimp.NewMockClient(&mock)

// perform some actions/operations here

if mock.PostCalls != 1 {
    t.Errorf("expected 1 PostCall, got %d", mock.PostCalls)
}
```

### Mocking the entire client
Furthermore, to mock the entire `mailchimp.Client` value, you can use the `mailchimp.ClientMock` which satisfies the interface for a regular `mailchimp.Client` but it leaves the application developer to define its behaviour. Any of the receiver functions can be mocked, and all of them are paired with a counter that signifies the amount of times the function has been called thus far. The usage of `mailchimp.ClientMock` is very similar to that of `mailchimp.MailChimpProviderMock` and an example is showcased below.

```go
mock := mailchimp.ClientMock{
    PingMock: func() error {
        return errors.New("could not connect to MailChimp")
    }
}

// inject the dependency and perform testing

if mock.PingCalls != 1 {
    t.Errorf("expected 1 call to Ping, but got %d", mock.PingCalls)
}
```

The `mailchimp.ClientMock` struct is shown below.

```go
type ClientMock struct {
	PingMock  func() error
	PingCalls int

	CreateListMock  func(List) (List, error)
	CreateListCalls int
	FetchListsMock  func() ([]List, error)
	FetchListsCalls int
	FetchListMock   func(string) (List, error)
	FetchListCalls  int
	UpdateListMock  func(string, List) (List, error)
	UpdateListCalls int
	DeleteListMock  func(string) error
	DeleteListCalls int

	BatchMock            func(string, []Member) error
	BatchCalls           int
	BatchWithUpdateMock  func(string, []Member) error
	BatchWithUpdateCalls int

	FetchMemberTagsMock       func(string, string) ([]Tag, error)
	FetchMemberTagsCalls      int
	UpdateMemberTagsMock      func(string, string, []Tag) error
	UpdateMemberTagsCalls     int
	UpdateMemberTagsSyncMock  func(string, string, []Tag) error
	UpdateMemberTagsSyncCalls int

	ArchiveMemberMock  func(string, string) error
	ArchiveMemberCalls int

	CreateWebhookMock  func(webhook Webhook) (Webhook, error)
	CreateWebhookCalls int
	FetchWebhookMock   func(listID string, webhookID string) (Webhook, error)
	FetchWebhookCalls  int
	DeleteWebhookMock  func(listID string, webhookID string) error
	DeleteWebhookCalls int
}
```
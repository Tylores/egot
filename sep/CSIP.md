# Common Smart Inverter Profile (CSIP) Test Procedures

The test title includes an indication of which category for which the test is required: 
- Device Client **[C]**
- Aggregator Client **[A]**
- Server **[S]**
- Test **[T]**

In the test descriptions, Device Clients and Aggregator Clients are both represented with the designation as a Client **[C]**

## Core Fuction Test Procedures

### 1. CORE-001 - HTTP Request [S]

**Procedure:**

1. **[T]** Record the Client/Server communications.
2. **[TC]** Retrieve the `DeviceCapability` resource from the Server using the supported HTTP and IP address.
3. **[TC]** Retrieve the `DeviceCapability` resource from the Server using the supported HTTPS and IP address.
4. **[TC]** Process the retrieved `DeviceCapability` resource and select an item in the `DeviceCapability` resource. For example, `DERProgram`.
5. **[TC]** Based on the WADL definition for the resource as defined in the IEEE 2030.5 standard, do GET, PUT, POST, DELETE operation as permitted using correct HTTP headers and payload (if required).
6. **[TC]** Based on the WADL definition for the resource as defined in the IEEE 2030.5 standard, do GET, PUT, POST, DELETE operation as disallowed using correct HTTP headers and payload (if required). The Server shall correctly return a 405 Method Not allowed HTTP response for the invalid operation.
7. **[TC]** Based on the WADL definition for the resource as defined in the IEEE 2030.5 standard, do GET operation as permitted with incorrect HTTP header. The Server shall correctly return a 400 Bad Request HTTP response for the invalid operation.
8. **[TC]** Do an invalid method operation, such as HTTP FOO, on the same resource. The Server shall correctly return a 501 Not Implemented HTTP response for the invalid method used.
9. Repeat this test by iterating through all subordinate resources from the parent resource presented in the `DeviceCapability` payload.

### 2. CORE-002 - HTTP Response [S]

**Procedure:**

1. **[T]** Record the Client/Server communications.
2. **[TC]** Retrieve the `DeviceCapability` resource from the Server using the supported HTTP and IP address, which shall return a 200 OK response from the server.
3. **[TC]** Select a resource that supports a POST or PUT. Next, perform an HTTP POST or PUT operation on that resource which should trigger a 201 Created or 204 No Content response.
4. **[TC]** Select a resource that supports a POST or PUT. Next, perform an HTTP GET operation on that resource with missing Host header in the HTTP header, which should trigger a 400 Bad Request response.
5. **[TC]** Select any resource that supports GET. Next, perform an HTTP GET operation on that resource with invalid URI, which should trigger a 404 Not Found response.
6. **[TC]** Select any resource that does not support POST. Next, perform an HTTP POST operation on that resource which should trigger a 405 Method Not Allowed response with Allow header with a list of valid methods for that requested resource.
7. **[TC]** Select any resource that supports GET. Next, perform an HTTP FOO on that resource which should trigger a 501 Not Implemented response.
8. Repeat Steps 2 through 7 on resources subordinate to the parent resource included in the `DeviceCapability` payload.

### 3. CORE-003 - Polling Interaction [C, A]

**Procedure:**

1. **[T]** Record the Client/Server communications.
2. **[C]** Retrieve the `DeviceCapability` resource from the Server using the supported HTTP and IP address and find the `EndDeviceListLink` element.
3. **[C]** Perform an HTTP GET operation on the `EndDeviceListLink` URI and search through the `EndDeviceList` payload to find if an `EndDevice` instance is included that matches the identity of the Client device. For example, SFDI/LFDI.
4. **[C]** Process the `EndDevice` instance returned by the Server and perform an HTTP GET operation on the `RegistrationLink` to retrieve the `Registration: PIN` assigned to the Client. Confirm the PIN is 111115.
5. **[C]** Process the `EndDevice` instance returned by the Server and perform an HTTP GET operation on the `FunctionSetAssignmentsListLink` to retrieve the various FSA assigned resources assigned to the Client.
6. **[C]** Using the `FunctionSetAssignments` instance, perform an HTTP GET operation on the `DERProgramListLink` to retrieve the `DERProgramList` assigned to the Client.
7. **[C]** Parse the returned `DERProgramList` payload to find how many `DERProgram` instances are included.
8. **[C]** Set the `pollRate` attribute in each of `DeviceCapability`, `EndDeviceList`, `FunctionSetAssignmentsList`, `DERProgramList` to verify the ability to observe the specified poll rate.

### 4. CORE-004 - List Handling [S]

**Procedure:**

1. **[T]** Record the Client/Server communications.
2. **[TC]** Retrieve the `DeviceCapability` resource from the Server using the supported HTTP and IP address.
3. **[TC]** Process the retrieved `DeviceCapability` resource and perform an HTTP GET operation on the `EndDeviceList` (or `FunctionSetAssignmentsList` or `DERProgramList`) resource link included in the `DeviceCapability` with query string parameter. For example, HTTP GET `http://ipaddress/EndDeviceList?s=0`. This query string search shall return the first element in the list.
4. **[TC]** Perform an HTTP GET operation on the list resource link with no query string parameter. For example, HTTP GET `http://ipaddress/EndDeviceList`. This request shall return the first element in the list.
5. **[TC]** Perform an HTTP GET operation with time-based query string parameter. For example, `http://ipaddress/EndDeviceList?a=time`, where time is current time expressed as a valid `TimeType`. This request shall return the first element in the list (if not time-based) or the element after the specified time.
6. **[TC]** Perform an HTTP GET operation with query string parameter: `http://ipaddress/EndDeviceList?l=0`. This request shall return an empty list.
7. **[TC]** Perform an HTTP GET operation with query string parameter: `http://ipaddress/EndDeviceList?1=1000`. This request shall return up to 1000 elements in the list depending on the number of elements that exist.
8. **[TC]** Perform an HTTP GET operation with time-based query string parameter: `http://ipaddress/EndDeviceList?s=3&a=time`. This request shall return the contents of the ordinal elements in the list beginning s=3 since none of these are time based.
9. **[TC]** Perform an HTTP GET operation with query string parameter: `http://ipaddress/EndDeviceList?s=1&l=1&s=2`. This request shall return the contents beginning with the third and ignoring the second, s=2.
10. **[C]** Perform an HTTP GET operation with query string parameter: `http://ipaddress/EndDeviceList?l=0&1=2`. This request shall return an empty list.
11. **[TC]** Perform an HTTP GET operation with time-based query string parameter: `http://ipaddress/EndDeviceList?s=1&a=time&a=timeplus-two4 hours`. This request shall return contents of the ordinal elements in the list beginning with the S=1.
12. **[TC]** Perform an HTTP GET operation with query string parameter: `http://ipaddress/EndDeviceList?1=1&b=2`. This request shall return contents of the list starting with the first element and only that element.

### 5. CORE-005 - Basic Time [C, A, S]

**Procedure:**

1. **[T]** Record the Client/Server communications.
2. **[C]** Retrieve the `DeviceCapability` resource from the Server using the supported HTTP and IP address.
3. **[C]** Process the retrieved `DeviceCapability` resource and verify the Time resource available from the Server and included as a link.
4. **[S]** Confirm the quality metric of the Time resource is 7.
5. **[C]** Retrieve the Time resource from the Server by using the href information included in the `DeviceCapability` time resource.
6. **[C]** Process the retrieved Time resource, verify its quality metric to be 7 (intentionally uncoordinated) and synchronize the Client time using the information and display the updated time information on the Client display, if supported.

### 6. CORE-006 - Advanced Time [C, A, S]

**Procedure:**

1. **[T]** Record the Client/Server communications.
2. **[C]** Retrieve the `DeviceCapability` resource from the Server using the supported HTTP and IP address.
3. **[C]** Process the retrieved `DeviceCapability` resource and verify the Time resource available from the Server and included as a link.
4. **[C]** Retrieve the Time resource from the Server by using the href information included in the `DeviceCapability` Time resource.
5. **[S]** Change the time forward by 1 hour which should trigger a `LogEvent` indicating `TM_TIME_ADJUSTED`.
6. **[C]** Retrieve the Time resource from the Server by using the href information included in the `DeviceCapability` Time resource to check the updated time.
7. **[C]** Retrieve the `LogEventList` from the Server `SelfDevice` and verify a `TM_TIME_ADJUSTED` `LogEvent` was generated.

### 7. CORE-009 - Advanced End Device [C, A, S]

**Procedure:**

1. **[T]** Record the Client/Server communications.
2. **[C]** Retrieve the `DeviceCapability` resource from the Server using the supported HTTP and IP address and find the `EndDeviceListLink` element.
3. **[C]** Perform an HTTP GET operation on the `EndDeviceListLink` URI and search through the `EndDeviceList` payload to find if an `EndDevice` instance is included that matches the identity of the Client device. If found, skip next step.
4. **[C]** Using the Location of the created `EndDevice` instance returned by the Server, perform an HTTP GET operation on that Location. On successful GET operation, perform an HTTP GET operation on the `RegistrationLink` href to find the PIN value for the Client device.
5. **[C]** Process the returned `Registration` resource and search for the PIN element and find its value. Verify the PIN value is the same PIN value the Client device has preregistered.
6. **[C]** Using the `EndDevice` instance resource returned in step 4, find the `DERListLink` information. Do an HTTP PUT on the `DERListLink` using the href attribute and updated values for `DERCapabilities`, `DERSettings`, `DERStatus` or `DERAvailability`.

### 8. CORE-010 - Function Set Assignments [C, A, S]

**Procedure:**

1. **[T]** Record the Client/Server communications.
2. **[S]** Configure resources as outlined in Setup and wait for Client to acquire resources.
3. Wait two minutes after resources before terminating the test.
4. **[C]** Retrieve the `DeviceCapability` resource from the Server using the supported HTTP and IP address and find the `EndDeviceListLink` element.
5. **[C]** Perform an HTTP GET operation on the `EndDeviceListLink` URI and search through the `EndDeviceList` payload to find if an `EndDevice` instance is included that matches the identity of the Client device.
6. **[C]** On successful identity match, perform an HTTP GET operation on the `RegistrationLink` href to find the PIN value for the Client device.
7. **[C]** Process the returned Registration resource and search for the PIN element and find its value. Verify the PIN value is the same PIN value the Client device has preregistered.
8. **[C]** Using the `EndDevice` instance resource returned in step 4, find the `FunctionSetAssignmentsListLink` information.
9. **[C]** Starting with the `FunctionSetAssignmentsListLink`, follow the link to find (1) the `FunctionSetAssignmentsList`, (2) the `FunctionSetAssignments` instance, (3) the `DERProgramListLink`, (4) the `DERProgramList`, and (5) the `DERPrograms`.
10. **[C]** GETs all the `DERPrograms`, or groups, assigned to it.

### 9. CORE-011 - Advanced Function Set Assignments [C, A, S]

**Procedure:**

1. **[T]** Record the Client/Server communications.
2. **[C]** Retrieve the `DeviceCapability` resource from the Server using the supported HTTP and IP address and find the `EndDeviceListLink` element.
3. **[C]** Perform an HTTP GET operation on the `EndDeviceListLink` URI and search through the `EndDeviceList` payload to find if an `EndDevice` instance is included that matches the identity of the Client device. If found, skip next step.
4. **[C]** Using the Location of the created `EndDevice` instance returned by the Server, perform an HTTP GET operation on that Location. On successful GET operation, the `EndDevice` instance payload returned by the Server shall include relevant subordinate resources assigned to the Client.
5. **[C]** Process the `EndDevice` instance returned by the Server and perform an HTTP GET operation on the `RegistrationLink` to retrieve the `Registration: PIN` assigned to the Client. Confirm the PIN is 111115.
6. **[C]** Process the `EndDevice` instance returned by the Server and perform an HTTP GET operation on the `FunctionSetAssignmentsListLink` to retrieve the various FSA assigned resources assigned to the Client.
7. **[C]** If there is more than one `FunctionSetAssignments` in the `FunctionSetAssignementsList` payload, do additional HTTP GET on the rest of the `FunctionSetAssignments` and sort it based on the Primary Key (mRID) to find the priority ordering.
8. **[C]** Using the highest priority `FunctionSetAssignments` instance, perform an HTTP GET operation on the `DERProgramListLink` to retrieve the various `DERProgramList` assigned to the Client. If no `DERProgramListLink` is found, go back to step 7 and select the next higher priority `FunctionSetAssignments` instance.
9. **[C]** Parse the returned `DERProgramList` payload to find how many `DERProgram` instances are included. If there is more than one, iterate through all `DERProgram` instances and apply the `DERProgram` list priority processing to select the highest priority `DERProgram`.
10. **[C]** Perform an HTTP GET operation on the `DERProgram` selected and process the returned payload. Process the elements in the `DERProgram` payload and do HTTP GET operations to get subordinate resources. When the DER event is found through `DERControl`, schedule it on the Client.

### 10. CORE-012 - Basic DER Program/Control [C, A, S]

**Procedure:**

1. **[T]** Record the Client/Server communications.
2. **[C]** Process the received `DERProgramList` assigned to the Client through the FSA instance and find the highest priority `DERProgram`/`DERControl` to process by following the list ordering rules. The DER-Client shall schedule the selected `DERProgram`/`DERControl`.
3. **[C]** After the correct `DERProgram`/`DERControl` is selected, Client shall issue subsequent HTTP GET requests to retrieve the subordinate resources associated with the `DERProgram` and `DERControl`, such as `DERCurveListLink` with one or more `DERCurveList` instances.
4. **[C]** For the `DERCurve` resources in the `DERCurveListLink`, the Client shall do HTTP GET requests to retrieve the 10 curve points and their values. The Client shall use these to construct the DER curve-based control, which shall be communicated to the internal inverter system and the schedule information.
5. **[C]** Until the selected `DERProgram`/`DERControl` becomes active, the Client shall use the `DefaultDERControl` value by setting the DER control mode using the `DefaultDERControl` attribute.
6. **[C]** Periodically poll the associated `DERControl` resource to check the Status of the DER event. Client completes the DER event either based on the `DERControl` status value and/or start time plus duration interval.

### 11. CORE-013 - Advanced DER Program/Control [C, A, S]

**Procedure:**

1. **[T]** Record the Client/Server communications.
2. **[C]** Process the received `DERProgramList` assigned to the Client through the FSA instance and find the highest priority `DERProgram`/`DERControl` to process by following the list ordering rules. The DER-Client shall schedule the selected `DERProgram`/`DERControl`.
3. **[C]** Client shall issue subsequent HTTP GET requests to retrieve the subordinate resources associated with each of the `DERProgram` and `DERControl`.
4. **[C]** Until the selected `DERProgram`/`DERControl` becomes active, the Client shall use the `DefaultDERControl` value by setting the DER control mode using the `DefaultDERControl` attribute.
5. **[C]** Repeat steps 4 & 5 for the remaining `DERProgram`/`DERControl` instances assigned to the Client device.

### 12. CORE-014 - Basic DER Settings (Power Generating) [C, A, S]

**Procedure:**

1. **[T]** Record the Client/Server communications.
2. **[C]** Perform an HTTP GET operation on the `DERList` resource assigned through its pre-registered `EndDevice` instance on the Server.
3. **[C]** Do an HTTP PUT to the DER resources with the list of links (and current information populated), which it maintains on the `DERList` resource presented through its `EndDevice` instance.
4. **[C]** Perform an HTTP GET operation on the `DERCapability` payload for this `EndDevice` to check the type element to find if the device has support for power generation.
5. **[C]** Perform an HTTP GET operation on the `DERSettings` payload for this `EndDevice`. This `DERSettings` should reflect the current settings of the individual DER Client device.
6. **[C]** Based on the `DERCapability` and `DERSettings` payloads, find presence of mandatory items for power generation DER devices (e.g. `opModMaxLimW`, `setMaxW`, `setMaxVA`, etc).
7. **[C]** Reactive Power (VAr) Test: If the `DERCapability/modesSupported` indicates support for `opModFixedVAr`, do the next test step. Otherwise, skip to Power Factor test section.
8. **[C]** Based on the `DERCapability` and `DERSettings` payloads, find presence of mandatory items for Reactive Power capable DER devices.
9. **[C]** Power Factor (PF) Test: If the `DERCapability/modesSupported` indicates support for `opModFixedPFInjectw`, do the next test step.
10. **[C]** Based on the `DERCapability` and `DERSettings` payloads, find presence of mandatory items for Power Factor capable DER devices.

### 13. CORE-018 - Basic Subscription [A, S]

**Procedure:**

1. **[T]** Record the Client/Server communications.
2. **[C]** Search the retrieved list of `EndDevice` instances from the Server for the Client SFDI/LFDI value. Next, find the `SubscriptionListLink` elements to confirm the Server supports subscription.
3. **[C]** Search the `FSAList` assigned to the Client `EndDevice` and verify it can be subscribed to.
4. **[C]** Prepare a subscription request that includes the `FSAList` and send an HTTP POST request.
5. **[S]** For the received subscription request from the Client, process and send an HTTP 201 Created or HTTP 204 No Content response with the correct HTTP Location header.
6. **[S]** Cause a change for one the elements of the FSA resource for the Client, which shall cause a notification message. At such resource change, send a notification message using the URI information.
7. **[C]** When the Notification message is received from the Server, process the incoming Notification message and payload, including validation.
8. **[C]** Perform an HTTP GET request on the href of the resource included in the Server Notification message. Process the payload returned from the Server from the HTTP GET request and compare to the Notification body payload (they shall be identical).

### 14. CORE-019 - Advanced Subscription [A, S]

**Procedure:**

1. **[T]** Record the Client/Server communications.
2. **[C, S]** Search the retrieved list of `EndDevice` instances from the Server for the Client SFDI/LFDI value. Find the `SubscriptionListLink` elements and confirm this `EndDevice` instance itself can be subscribed to.
3. **[C]** Prepare a subscription request for this `EndDevice` instance with required elements/attributes and send an HTTP POST request to the Server.
4. **[S]** Receive the sent subscription request from step 3 and process the request and respond with an HTTP 201 Created or HTTP 204 No Content response.
5. **[C]** Search the `FSAList` assigned to the Client `EndDevice` and traverse through the found FSA instance for an element that can be subscribed to.
6. **[C]** For the FSA, prepare a subscription request with all required elements/attributes and send an HTTP POST request.
7. **[S]** For the received subscription request from the Client, process and send an HTTP 201 Created or HTTP 204 No Content response.
8. **[S]** Cause a change for one the elements of the FSA resource for the Client. Send a notification message using the URI information sent by the Client and include the updated FSA resource payload.
9. **[C]** When the Notification message is received from the Server, process the incoming Notification message and payload, and respond back with HTTP 201 Created or HTTP 204 No Content message.
10. **[S]** The resource change from step 8 shall also cause a related Notification message for the `EndDevice` subscribed to by the Client. The Server shall send a notification message using the URI information.
11. **[C]** When the Notification message is received from the Server for the `EndDevice` resource update, process it and respond back.
12. **[C]** Perform an HTTP GET request on the href of the resource included in the Server notification message. Process the payload returned from the Server and verify that it is identical to the notification body payload. Repeat for each notification received.
13. **[S]** Cancel the outstanding FSA subscription by sending a notification message to Client indicating `notification status=1` (Subscription canceled).
14. **[C]** Receive the Notification message that indicates cancellation sent by the Server, process the message, and update the internal subscription state.

### 15. CORE-021 - Randomized Events [C, A, S]

**Procedure:**

1. **[T]** Record the Client/Server communications.
2. **[C]** Process the received `DERProgramList` assigned to the Client through the FSA instance and find the highest priority `DERProgram`/`DERControl`. The DER-Client shall schedule the selected `DERProgram`/`DERControl` with the correct `randomizeStart` and `randomizeDuration` values applied.
3. **[C]** Client shall issue subsequent HTTP GET requests to retrieve the subordinate resources associated with all of the `DERProgram` and `DERControl` instances, such as `DERCurveListLink`.
4. **[C]** For the `DERCurve` resources in the `DERCurveList`, the Client shall do HTTP GET requests to retrieve the 10 curve points and their values. The Client shall use these to construct the DER curve-based control (e.g. `start time + duration + randomizationStart + randomizeDuration`).
5. **[C]** Until the selected `DERProgram`/`DERControl` becomes active, the Client shall use the `DefaultDERControl` value.
6. **[C]** Periodically poll the associated `DERControl` resource to check the Status of the DER event. Client completes the DER event based on the calculated randomized period.
7. **[C]** Repeat steps #4 and 5 for rest of the `DERProgram`/`DERControl` instances that are scheduled.

### 16. CORE-022 - Responses [C, A, S]

**Procedure:**

1. **[T]** Record the Client/Server communications.
2. **[C]** Process the received `DERProgramList` assigned to the Client through the FSA instance and find the highest priority `DERProgram`/`DERControl` to process by following the list ordering rules. The DER-Client shall schedule the selected `DERProgram`/`DERControl`.
3. **[C]** After the correct `DERProgram`/`DERControl` is selected, Client shall issue subsequent HTTP GET requests to retrieve the subordinate resources associated with the `DERProgram` and `DERControl`.
4. **[C]** For the `DERCurve` resources in the `DERCurveListLink`, the Client shall do HTTP GET requests to retrieve the 10 curve points and their values and schedule the event.
5. **[C]** Until the selected `DERProgram`/`DERControl` becomes active, the Client shall use the `DefaultDERControl` value by setting the DER control mode using the `DefaultDERControl` attribute.
6. **[S]** Update the `DERControl#1` `currentStatus`/`dateTime` values at each state of the DER event by following its event schedule. Before the event completes its schedule, cancel the event by updating its `currentStatus` value to 2 (canceled).
7. **[C]** Periodically poll the associated `DERControl#1` resource to check the Status of the DER event and send the requested response at the `responseRequired` attribute, such as a Response message with status=1 (Received) and status=2 (Started). Because this event was canceled in step 5, the Client shall notice such cancellation and send a response with status=6 (cancellation).
8. **[S]** Update the `DERControl#2` `currentStatus`/`dateTime` values at each state of the DER event by following its event schedule.
9. **[C]** Periodically poll the associated `DERControl#2` resource to check the Status of the DER event and send the requested response at the `responseRequired` attribute. Client completes the DER event either based on the `DERControl#2` Status value and/or start time plus duration interval.
10. **[S]** Update the `DERControl#3` `currentStatus`/`dateTime` values at each state of the DER event by following its event schedule.
11. **[C]** Periodically poll the associated `DERControl#3` resource to check the Status of the DER event and send the requested response at the `responseRequired` attribute. Client completes the DER event either based on the `DERControl#3` Status value and/or start time plus duration interval.

## Basic Function Test Procedures

Here are the extracted procedures for the Basic Function Tests found in the provided document:

### 1. BASIC-001 - DER Identification [C, A, S]

**Procedure:**

1. **[T]** Record the Client/Server communications.
2. **[C]** Retrieve the `DeviceCapability` resource from the Server using the supported HTTP and IP address and find the `EndDeviceListLink` element.
3. **[C]** Perform an HTTP GET operation on the `EndDeviceListLink` URI and search through the `EndDeviceList` payload to find if an `EndDevice` instance is included that matches the identity of the Client device.
4. **[C]** Using the `EndDevice` instance found from the previous step, validate the SFDI and LFDI values.
5. **[C]** Using the `RegistrationLink`/PIN instance found from the `EndDevice` instance, validate the PIN values.
6. **[C]** For an aggregator client, iterate through rest of the `EndDevice` instances included in the original `EndDeviceList` from the previous test step and validate each set of SFDI/LFDI and PIN values.
7. **[C]** Using the `FunctionSetAssignmentsListLink` from test step 4, perform an HTTP GET operation on the included resource inside the `FunctionSetAssignmentsList` the link refers to.

### 2. BASIC-002 - Basic Group Management [C, A, S]

**Procedure:**

1. **[T]** Record the Client/Server communications.
2. The Client shall have retrieved all of its seven `FunctionSetAssignments` that correspond to the grid topology-based diagram to continue execution of this test.
3. **[C]** Using the FSA instances retrieved from the CORE-010 test, process the elements in the FSA instance to do additional HTTP GET requests for all of its subordinate resources. Validate there is a `Time` resource assigned to the FSA instance.
4. **[C]** Do HTTP GET requests on the `DERProgramListLink` and `TimeLink` hrefs for the FSA instance using `l=255` query string search parameter.
5. **[C]** Process the information in the returned `DERProgramList` from the previous step and find which `DERProgram` has the highest priority based on IEEE 2030.5 event priority determination rules.
6. **[C]** Process the contents of the highest priority `DERProgram` and find which `DERControl` to make active: `DefaultDERControl` or `DERControlListLink` associated `DERControl`. Do additional HTTP GET requests for subordinate resources to find which `DERControl` to make active.
7. **[C]** Based on the determination in the previous step, make the appropriate `DERControl` active (based on the test setup, it shall be the `DefaultDERControl` that becomes active).
8. **[C]** Periodically poll the FSA resource for updated information in the FSA instance.

### 3. BASIC-003 - Advanced Group Management [C, A, S]

**Procedure:**

1. **[T]** Record the Client/Server communications.
2. The Client shall have retrieved all of its seven `FunctionSetAssignments` that correspond to the topology specified in Figure 3 to continue execution of this test.
3. **[S]** Before the scheduled `DERProgram`/`DERControl` starts, change the `FSAList` assigned to this Client by removing the specific feeder in the diagram and adding a new feeder to the topology.
4. **[C]** Poll the FSA resource for updated information in the FSA instance, which shall cause the Client to notice the updated `FSAList` and issue subsequent HTTP GETs to retrieve the updated subordinate information.
5. **[C]** Based on the updated resources from the previous step, schedule a different `DefaultDERControl` included in the updated `FSAList` `DERProgram`.

### 4. BASIC-004 to BASIC-014 - Basic Inverter Controls

*(Note: The procedures for BASIC-004 through BASIC-014 follow a uniform procedural structure applied to different control functions such as Low/High Voltage Ride-Through, Low/High Frequency Ride-Through, Volt/Var, Ramp Rates, Fixed Power Factor, Connect/Disconnect, Limit Max Active Power Mode, Volt-Watt, Frequency-Watt, and Set Active Power Modes.)*

**Procedure:**

1. **[T]** Record the Client/Server communications.
2. **[C]** The Client shall have retrieved all its FSA group and associated resources, including all subordinate resources for `DERProgram`, `DERControl`, and `DefaultDERControl`, `DERCurveList` and others. The Client would have processed all resource information for the FSA to find the priorities of each `DERProgram` and associated `DERControl` resources.
3. **[C]** Using the `DefaultDERControl` information associated with the highest priority `DERProgram`, the Client shall activate the `DefaultDERControl` until the Service Point DER event becomes active.
4. **[C]** At the effective start time of the highest priority `DERProgram`/`DERControl`, it shall check the `currentStatus` of the `DERControl` by executing an HTTP GET on the selected `DERControl` href. After the `currentStatus` field is verified, Client shall activate this `DERControl` using the subordinate resource information. If a response is required, Client shall send a series of response messages based on the response requirements specified.
5. **[S]** Process the sent response messages and verify they reflect the `currentStatus` of the actual DER event.
6. **[C]** Complete the current `DERControl` event at the effective end time and send the required response message to the Server. If there is no other active event, it defaults to the `DefaultDERControl` information associated with the highest priority `DERProgram`.
7. **[S]** Verify all the immediate or curve-based control is exercised in this test for the Client by checking the response message sent by the Client for each `DERControl` event.

### 5. BASIC-016 - Event - 2 DERP, 2 DDERC, 0 DERC [C, A, S]

**Procedure:**

1. **[T]** Record the Client/Server communications.
2. **[C]** GETs the SY `DERProgramList`.
3. **[C]** GETs the SY `DERProgram`.
4. **[C]** GETs the SY `DefaultDERControl`.
5. **[C]** GETs the SP `DERProgramList`.
6. **[C]** GETs the SP `DERProgram`.
7. **[C]** GETs the SP `DefaultDERControl`.
8. **[C]** Applies the SP `DefaultDERControl` because it has a higher priority (lower primacy value) than the SY `DefaultDERControl`.

### 6. BASIC-017 - Event - 1 DERP, 0 DDERC, 1 DERC [C, A, S]

**Procedure:**

1. **[T]** Record the Client/Server communications.
2. **[C]** GETs the SY `DERProgramList`.
3. **[C]** GETs the SY `DERProgram`.
4. **[C]** Periodically GETs the SY `DERControlList`.
5. **[C]** GETs of the created SY `DERControl`.
6. **[C]** POSTs response with status 1 (Event Received) to the Server.
7. **[C]** Applies the SY `DERControl` at the correct start time for the correct duration.
8. **[C]** POSTs response with status 2 (Event Started) at start time of the event.
9. **[C]** POSTs response with status 3 (Event Completed) after duration has elapsed.

### 7. BASIC-018 - Event - 1 DERP, 1 DDERC, 1 DERC [C, A, S]

**Procedure:**

1. **[T]** Record the Client/Server communications.
2. **[C]** GETs the SY `DERProgramList`.
3. **[C]** GETs the SY `DERProgram`.
4. **[C]** GETs the SY `DefaultDERControl`.
5. **[C]** Applies the SY `DefaultDERControl` because there is no active event.
6. **[C]** Periodically GETs the SY `DERControlList`.
7. **[C]** GETs the created SY `DERControl`.
8. **[C]** POSTs response with status 1 (Event Received) to the Server.
9. **[C]** Applies the SY `DERControl` at the correct start time for the correct duration.
10. **[C]** POSTs response with status 2 (Event Started) at start time.
11. **[C]** POSTs response with status 3 (Event Completed) after duration has elapsed.
12. **[C]** Applies the SY `DefaultDERControl` after the completion of the event.

### 8. BASIC-019 - Event - 1 DERP, 1 DDERC, 2 Non-overlapping Similar DERC [C, A, S]

**Procedure:**

1. **[T]** Record the Client/Server communications.
2. **[C]** GETs the SY `DERProgramList`.
3. **[C]** GETs the SY `DERProgram`.
4. **[C]** GETs the SY `DefaultDERControl`.
5. **[C]** Applies the SY `DefaultDERControl` because there is no active event.
6. **[C]** Periodically GETs the SY `DERControlList`.
7. **[C]** GETs of the SY `DERControl` event created in setup step 2.
8. **[C]** POSTs response with status 1 (Event Received) for the event created in setup step 2 to the Server.
9. **[C]** Applies the SY `DERControl` created in Setup Step 2 at the correct start time for the correct duration.
10. **[C]** POSTs response with status 2 (Event Started) at start time of event created in setup step 2.
11. **[C]** POSTs response with status 3 (Event Completed) after duration of the event created in setup step 2 has elapsed.
12. **[C]** Applies the SY `DefaultDERControl` after the completion of the event created in setup step 2.
13. **[C]** GETs of the SY `DERControl` event created in setup step 3.
14. **[C]** POSTs response with status 1 (Event Received) for the event created in setup step 3 to the Server.
15. **[C]** Applies the SY `DERControl` created in setup step 3 at the correct start time for the correct duration.
16. **[C]** POSTs response with status 2 (Event Started) at start time of event created in setup step 3.

### 9. BASIC-020 - Event - 2 DERP, 2 DDERC, 2 Non-overlapping Similar DERC [C, A, S]

**Procedure:**

1. **[T]** Record the Client/Server communications.
2. **[C]** GETs the SY `DERProgramList`.
3. **[C]** GETs the SY `DERProgram`.
4. **[C]** GETs the SY `DefaultDERControl`.
5. **[C]** GETs the SP `DERProgramList`.
6. **[C]** GETs the SP `DERProgram`.
7. **[C]** GETs the SP `DefaultDERControl`.
8. **[C]** Applies the SP `DefaultDERControl` because there is no active event and has a higher priority than the SY `DefaultDERControl`.
9. **[C]** Periodically GETs the SY `DERControlList`.
10. **[C]** Periodically GETs the SP `DERControlList`.
11. **[C]** GETs of the created SY `DERControl`.
12. **[C]** POSTs response with status 1 (Event Received) to the Server.
13. **[C]** Applies the SY `DERControl` at the correct start time for the correct duration.
14. **[C]** POSTs response with status 2 (Event Started) for the SY event at start time.
15. **[C]** POSTs response with status 3 (Event Completed) for the SY event after duration has elapsed.
16. **[C]** Applies the SP `DefaultDERControl` after the completion of the event.
17. **[C]** GETs of the SP `DERControl`.
18. **[C]** POSTs response with status 1 (Event Received) to the Server.
19. **[C]** Applies the SP `DERControl` at the correct start time for the correct duration. 
20. **[C]** POSTs response with status 2 (Event Started) for the SP event at start time.
21. **[C]** POSTs response with status 3 (Event Completed) for the SP event after duration has
elapsed.
22. **[C]** Applies the SP `DefaultDERControl` after the completion of the event because it
has a higher priority than the SY `DefaultDERControl`.

Here are the procedures for the remaining **Basic Function Tests** extracted and synthesized from the document patterns:

### 10. BASIC-021 - Event - 2 DERP, 2 DDERC, 2 Overlapping Similar DERC - System DERC followed by Service Point DERC before Start of System DERC [C, A, S]

**Procedure:**

1. **[T]** Record the Client/Server communications.
2. **[C]** GETs the SY `DERProgramList`.
3. **[C]** GETs the SY `DERProgram`.
4. **[C]** GETs the SY `DefaultDERControl`.
5. **[C]** GETs the SP `DERProgramList`.
6. **[C]** GETs the SP `DERProgram`.
7. **[C]** GETs the SP `DefaultDERControl`.
8. **[C]** Applies the SP `DefaultDERControl` because there is no active event and it has a higher priority than the SY `DefaultDERControl`.
9. **[C]** Periodically GETs the SY `DERControlList`.
10. **[C]** GETs the created SY `DERControl`.
11. **[C]** POSTs response with status 1 (Event Received) to the Server for the SY `DERControl`.
12. **[C]** Periodically GETs the SP `DERControlList`.
13. **[C]** GETs the SP `DERControl`.
14. **[C]** POSTs response with status 1 (Event Received) to the Server for the SP `DERControl`.
15. **[C]** Prior to the start time of the SY `DERControl`, POSTs response with status 7 (Event Superseded) to the Server for the SY `DERControl` because the SP event overlaps and has higher priority.
16. **[C]** Applies the SP `DERControl` at the correct start time for the correct duration.
17. **[C]** POSTs response with status 2 (Event Started) for the SP event at the correct start time.
18. **[C]** POSTs response with status 3 (Event Completed) for the SP event after duration has elapsed.
19. **[C]** Applies the SP `DefaultDERControl` after the completion of the event because it has a higher priority than the SY `DefaultDERControl`.

### 11. BASIC-022 - Event - 2 DERP, 2 DDERC, 2 Overlapping Similar DERC - Service Point DERC followed by System DERC [C, A, S]

**Procedure:**
*(Follows a similar sequence to BASIC-021, but tests the condition where the higher priority Service Point DERC is scheduled first. The lower priority System DERC is superseded immediately upon receipt or prior to execution because it overlaps with the already active or scheduled higher priority event.)*

1. **[T]** Record the Client/Server communications.
2. **[C]** GETs the associated `DERProgramList`, `DERProgram`, and `DefaultDERControl` for both SY and SP groups.
3. **[C]** Applies the SP `DefaultDERControl` (higher priority).
4. **[C]** Periodically GETs the SP and SY `DERControlList`.
5. **[C]** GETs the created SP `DERControl` and POSTs status 1 (Event Received).
6. **[C]** GETs the created SY `DERControl` and POSTs status 1 (Event Received).
7. **[C]** POSTs response with status 7 (Event Superseded) for the SY `DERControl` due to the overlapping higher-priority SP event.
8. **[C]** Applies the SP `DERControl` at its start time and POSTs status 2 (Event Started).
9. **[C]** POSTs status 3 (Event Completed) for the SP event when the duration elapses, then reverts to the SP `DefaultDERControl`.

### 12. BASIC-023 - Event - 2 DERP, 2 DDERC, 2 Overlapping Similar DERC - System DERC followed by Service Point DERC after Start of System Event [C, A, S]

**Procedure:**

1. **[T]** Record the Client/Server communications.
2. **[C]** GETs the SY `DERProgramList`, `DERProgram`, and `DefaultDERControl`.
3. **[C]** GETs the SP `DERProgramList`, `DERProgram`, and `DefaultDERControl`.
4. **[C]** Applies the SP `DefaultDERControl` because there is no active event and it has a higher priority.
5. **[C]** Periodically GETs the SY `DERControlList`.
6. **[C]** GETs the created SY `DERControl`.
7. **[C]** POSTs response with status 1 (Event Received) to the Server for the SY `DERControl`.
8. **[C]** Applies the SY `DERControl` at the correct start time.
9. **[C]** POSTs response with status 2 (Event Started) for the SY event at start time.
10. **[C]** Periodically GETs the SP `DERControlList`.
11. **[C]** GETs the SP `DERControl`.
12. **[C]** POSTs response with status 1 (Event Received) to the Server for the SP `DERControl`.
13. **[C]** At the start time of the SP `DERControl`, POSTs response with status 7 (Event Superseded) to the Server for the active SY `DERControl`.
14. **[C]** Applies the SP `DERControl` at the correct start time for the correct duration.
15. **[C]** POSTs response with status 2 (Event Started) for the SP event at correct start time.
16. **[C]** POSTs response with status 3 (Event Completed) for the SP event after duration has elapsed.
17. **[C]** Applies the SP `DefaultDERControl` after the completion of the event.

### 13. BASIC-024 - Event - 2 DERP, 2 DDERC, 2 Overlapping Independent DERC - System DERC followed by Service Point DERC before Start of System DERC [C, A, S]

**Procedure:**

1. **[T]** Record the Client/Server communications.
2. **[C]** GETs the SY `DERProgramList`, `DERProgram`, and `DefaultDERControl`.
3. **[C]** Applies the SY `DefaultDERControl`.
4. **[C]** GETs the SP `DERProgramList`, `DERProgram`, and `DefaultDERControl`.
5. **[C]** Applies the SP `DefaultDERControl`.
6. **[C]** Periodically GETs the SY `DERControlList` and SP `DERControlList`.
7. **[C]** GETs the SY `DERControl`.
8. **[C]** POSTs response with status 1 (Event Received) to the Server for the SY `DERControl`.
9. **[C]** Applies the SY `DERControl` at the correct start time for the correct duration.
10. **[C]** POSTs response with status 2 (Event Started) at start time for the SY `DERControl`.
11. **[C]** GETs the SP `DERControl`.
12. **[C]** POSTs response with status 1 (Event Received) to the Server for the SP `DERControl`.
13. **[C]** Applies the SP `DERControl` at the correct start time for the correct duration (these are independent controls, so both execute).
14. **[C]** POSTs response with status 2 (Event Started) at start time for the SP `DERControl`.
15. **[C]** POSTs response with status 3 (Event Completed) at the completion of the SY `DERControl`.
16. **[C]** Applies the SY `DefaultDERControl` after the completion of the event.
17. **[C]** POSTs response with status 3 (Event Completed) at the completion of the SP `DERControl`.
18. **[C]** Applies the SP `DefaultDERControl` after the completion of the event.

### 14. BASIC-025 & BASIC-026 - Event - 2 DERP, 2 DDERC, 2 Overlapping Independent DERC

**Procedure:**
*(These tests follow the exact same operational logic as BASIC-024, verifying that independent controls (e.g., one control manages active power while the other manages reactive power) can execute simultaneously without superseding one another, regardless of whether the Service Point or System event starts first.)*

1. **[T]** Record communications.
2. **[C]** Client retrieves lists, programs, and default controls for both levels.
3. **[C]** Client polls for new events and retrieves both the SP and SY `DERControls`.
4. **[C]** Client responds with Status 1 (Received) for both.
5. **[C]** Client executes the first event at its start time and responds with Status 2 (Started).
6. **[C]** Client executes the second independent event at its start time and responds with Status 2 (Started), running both controls concurrently.
7. **[C]** As each event's duration elapses, the Client responds with Status 3 (Completed) and reverts that specific parameter back to the corresponding `DefaultDERControl`.

### 15. BASIC-027 - Alarms [C, A, S]

**Procedure:**

1. **[T]** Record the Client/Server communications.
2. **[C]** The Client experiences a condition that triggers an alarm (e.g., hardware fault, communication loss).
3. **[C]** The Client generates a `LogEvent` formatted with the appropriate alarm code as specified by IEEE 2030.5.
4. **[C]** The Client performs an HTTP POST of the `LogEvent` to the `LogEventList` URI on the Server.
5. **[S]** The Server successfully receives the POST, parses the `LogEvent`, and returns a 201 Created response.
6. **[C]** When the alarm condition clears, the Client updates or posts a new `LogEvent` to clear the alarm state.

### 16. BASIC-028 - Inverter Status [C, A, S]

**Procedure:**

1. **[T]** Record the Client/Server communications.
2. **[C]** The Client periodically determines its operational status (e.g., operational state, connection status).
3. **[C]** The Client updates its local `DERStatus` resource.
4. **[C]** The Client performs an HTTP PUT operation to update the `DERStatus` resource on the Server using the associated `DERList` href.
5. **[S]** The Server receives the PUT request, updates the status, and returns a 204 No Content.
6. **[S]** Verify the updated `DERStatus` parameters on the server match the data transmitted by the Client.

### 17. BASIC-029 - Inverter Meter Reading [C, A, S]

**Procedure:**

1. **[T]** Record the Client/Server communications.
2. **[C]** The Client periodically measures its generation or consumption metrics.
3. **[C]** The Client updates its internal `MirrorMeterReading` resources.
4. **[C]** The Client posts the meter data by executing an HTTP POST to the Server's `MirrorMeterReading` or `DER` reading lists.
5. **[S]** The Server receives the meter reading, successfully validates the XML payload, and responds with a 201 Created.
6. **[S]** Verify the received readings conform to the required accuracy and required data points (e.g., Real Power, Apparent Power, Voltage) mandated by the CSIP specification.

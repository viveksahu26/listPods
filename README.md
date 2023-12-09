# listPods
## Write your own Custom Controller:
### But before that one question arises: Why to write own Custom Controllers?
    1) For me: motivation is to get feel of how Custom Controllers Work ?
    2) Writting something own will force you to research by own which results into self learning and growth.
    3) Satisfaction of building something by own.
    4) Learn by doing.

### Why Controllers and from where this term evolved ?
    I am not sure, but I think it evolved after the birth or revolution of Kubernetes. In Kubernetes, there is a concept of desired state and current state. 
    Desired State ---> What do you want.(Moreover it's like your future aim)
    Current State ---> Where are you now.(Where you stand right now.)

In real life desired state is your goal or say your future.  Whereas current state is your present state. Whereas in the world of Kuberenetes the desired state or current state is applied on resources. Resources such as Pods, Deployment, Service or custom resources.

So, In order to achieve desired state from current state, there should be someone to look out or watch out after resources progress and inform about the progress of that resource. And that's **where comes the use of Controllers**.

The job of Kubernetes Controllers is to watch out resources(or objects) current state with the help of Informers, and compare it with desired state provided by user through a `YAML` file. And the process of moving resource current state to desired state is done by the Controllers. It means Controllers has to continously fetch or pull information or data from Kuberentes(i.e. API Server) about current state of Kuberentes resources. But how it does that ?

### How does the Controller retrieves/fetch/pulls the resource information ?
To retrieve an resource or object information about current state of resources the controller has to sends request or reach out to a Kubernetes API Server. But the problem here occurs is, continously pinging to API Sever could lead to performance degradation, `why so`(will talk about it later). Let's take an example to understand it. We all know that post office duty is to deliver your package or letter whenever it comes. Now if a person start going to post offices everyday just to ask that is there any letter or post came on his name. Do you know post man will suggest you not need to come everyday, instead whenever any post(or any event) will come with your then we will reach out to you. Now, you would be thinking what is the problem if a man is comes to the post offics everyday. If one man comes, there is no problem with to the post office, but if this number is scaled or increases then it would lead to an problem or kind of disturbance to post office. Post office employee performance will be reduced because employee of post office will be busy on answering to each and every person whther there is any letter on their name, then how employee will do their work. Now get back to the controller and kubernets. Similarly, Kubernetes advices that: Hey, Controller we really appreciate your job/duty of comming to watch out resources but you don't need to come everytime. Instead if any changes occur in the resource which you are looking for then we will contact you. Finally Controllers aggreed and thanked to API server for saving his resources. Lastly Controller asked Kubernetes API Server that how it would be possible? In short who will be that guy who will approach me(controller) as any update being made to the resource(controller watching for), Kubernetes replied with a smile: `Informers` will do it for you.

### What is Informers ?
Informer is one which provides the information about resource to controllers as any changes occurs to the resource in Kuberentes cluster for which controller is looking for. So, Whenever there is any change in the status of the resource occurs internally then Informers triggers Controllers. `Trigger` is a technical term which means "letting you know if something happens". Similar to the post office guy who reach out to you whenever any letter comes with your name in the post office.

So, basically on Controller behalf, Informers retrives resource or objects information internally from Kubernetes and then inform to Controllers. Because of Informers Controllers job becomes quite easy now. Because of Controllers curiousity, it couldn't stopped himself from asking to Kubernetes that **how Informers does that ?**

### how Informers retrieves information internally from Kubernetes API Server ?
Internally, Informers queries the resource or object data and store it in a local cache. 
After storing data it compares with previous stored data(if any), if it detects any changes between those then it will update previously stored data with new one and then event will get triggered to tell about the changes in the state of resource to controllers. 

So, till now we know that why Informers were used: to save performance degrations which was caused by continous pinging to Kubenretes by controllers. And also unlike controllers, Informers don't have to ping to API Server becuase Informer is a internal part of Kubernetes whereas Controller is not.  Let's discover in more detail that:

### how Informers Works ?
Basically Informers has 3 components:
    1) Reflector
    2) DeltaFIFO
    3) Indexer

#### What is the role of Reflector ?
It watches specific resource or objects using K8s API. Resource may be in-built like Pods, Deployment, Service, etc and apart from in-built custom resources are alo there. So, reflector watches in-built as well as custom resources. And after watching resource status it adds events accordingly like `Added`, `Updated`, and `Deleted` to local cache i.e. `DeltaFIFO`.
`Added` represent to whether any resource added to the kubernetes cluster.
`Updated` represent to whether any resource updated to the kubernetes cluster.
`Deleted` represent to whether any resource deleted from the kubernetes cluster.

#### How Reflector Interact with K8S API ?
To talk to the k8s API. K8s API refers to the resource API.  It uses an `ListAndWatch` concept behind the scene. `ListAndWatch` = `List Resource` `+` `Watch Resource`. So, it lists all the resource corresponding to that k8s API and then Watch those resources. 

#### How ListAndWatch Works ?
##### How List Works ?
1) Firstly k8s API request is sent to get list of all those resource. In all resource there is one unique identity i.e. `ResourceVersion`, which represent the version or state of resource.  It's use comes in the Watch section. So, finally `List` helped us in retrieving all resources information with their  `ResourceVersion`.
##### How Watch Works ?
After succesfully fetching list of resources, the `Watch` will watch out `ResourceVersion`. If the version among newly fetched data and previously stored data remains same then there is no change in the resource state. But if their is change in the version of resource then resource is either Updated or Added, or Deleted. It compares previously stored `ResourceVersion` with newly fetched resources. If there is any difference in the `ResourceVersion` then it update the resource in the DeltaFIFO store. On the basis of comparision and depending on the type of change in the resource different events are added. If resource is newly Added, then `Added` events will be added in the DeltaFIFO, Similarly, if resource is updated then `Updated` events is added to DeltaFIFO store and similarly if resource is deleted then `Deleted` events is added in the DeltaFIFO store.

#### What is the role of DeltaFIFO ?
It stores resource events. The word `Delta` means change. This is special store that maintains `Delta`. So, finally from the reflector section with the help of ListAndWatch concept resource events are feeded or added in the DeltaFIFO store. But the question arises is:

#### Who will read these event stored in DeltaFIFO store ? or What is the use of Storing events in DeltaFIFO ?
It is used by the 3rd component of Informers i.e. `Cache Controller`. 

#### What is the role of Indexer ?


`NOTE`: Each and every resource or object has their own Controllers. 
    So, if resource exist then it's controller will also exist. And if controller exist then corresponding informers too exit.

### Till now what will learnt:
As we know each and every resource has their own controllers. And controllers helps in achieving from current state to desired state of the resource. The information about current state of resorce is provided by Informers in form of events triggering.

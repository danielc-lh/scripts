# capoeira

a lightweight library for choreographic programming via EPP-DI (Endpoint Projection as Dependency Injection) based on [ChoRus](https://github.com/lsd-ucsc/ChoRus). See the [paper](https://users.soe.ucsc.edu/~lkuper/papers/chorus-cp24.pdf) or the [docs](https://lsd-ucsc.github.io/ChoRus/introduction.html) for the core ideas.

# why? 
distributed code is hard to understand when you have multiple microservices interacting with each other.

it is not easy to see how data flows through your system. 

the goal of choreographic programming is to lay out the interaction between all parties in a global plan (a *choreography*) and then let each party run a program that performs the right actions at each step in the choreography. 

# what? 
a *choreography* is a global description of interactions between parties. 

it allows you to write a distributed protocol as a single program that deals with values or computations at multiple *locations*. 

using EPP (Endpoint Projection), the choreography can be projected out to each location as a simplified program that sends values to other locations, receieves values from other locations, or performs computations using values it currently has.

# examples
see [bookseller.go](./bookseller.go) for the classic bookseller protocol. 

a buyer + seller of books are communicating. the buyer has a title + a budget. the seller holds the inventory + list of corresponding prices.

1. the buyer sends the title of the book they want to the seller
2. the seller locally looks up the book in their inventory, returning the price if it is present and nil otherwise.
3. the buyer looks at the price and if is not nil, compares it against their budget. if it is within their budget, they send a messager to the seller to buy it. 
4. if the buyer wants to buy the book, the seller will respond to the buyer with the delivery date for the book.

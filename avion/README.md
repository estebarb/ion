Avion Web Framework
=======================

Avion is a experimental web framework over Ion.

Rationale behind Avion
-----------------

Avion is a experiment of creating a higher level web framework. Althought it
uses Ion, Avion is trying to implement a opinionated framework, capable of
rapid developing.

With that in mind, Avion will reuse the components of Ion, but will accelerate
the development of applications using several patterns, like:

- Controller based workflow: The basic building block in Avion is a controller,
  that is just a struct implementing the Avion Controller Interface. Those
  functions still are http.HandlerFunc, so it will retain the compatibility with
  net/http.
- RESTful routing: The controllers will implement a RESTful routing
  and actions. The routing can be direct, or nested.
- Easy way to define HTML views, or JSON views. The creation of API should be
  easy.

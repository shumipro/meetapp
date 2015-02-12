// es6
import Hello from './components/Hello.js'
var hello = new Hello('Hello')
hello.say()
hello.later().then(() => hello.say())

// react
import React from 'react'
import HelloReact from './jsx/Hello'

React.render(
  <HelloReact />,
  document.getElementById('example')
);
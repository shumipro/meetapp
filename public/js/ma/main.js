// es6 sample
import Hello from './components/Hello.js'
var hello = new Hello('Hello')
hello.say()
hello.later().then(() => hello.say())

// react sample
import React from 'react'
import HelloReact from './jsx/Hello'

React.render(
  <HelloReact />,
  document.getElementById('example')
);
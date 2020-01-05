import React, { useState } from 'react'
import {
  BrowserRouter as Router,
  Switch,
  Route,
  Link
} from 'react-router-dom'
import Home from './components/home'
import Index from './components/index'
import Package from './components/package'

function App() {

  

  return (
    <Router>
      <Switch>
        <Route path="/index">
          <Index />
        </Route>
        <Route path="/packages/:package">
          <Package />
        </Route>
        <Route path="/">
          <Home />
        </Route>
      </Switch>
    </Router>
  )
}

export default App

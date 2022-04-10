import React from 'react';
import logo from './logo.svg';
import './App.css';
import { LoginScreen } from './components/login'
import { GameScreen } from './components/game'

function App() {
  return (
    <GameScreen></GameScreen>
    // <LoginScreen></LoginScreen>
  );
}

export default App;

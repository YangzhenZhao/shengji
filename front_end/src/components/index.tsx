import React, { useState } from "react";
import { LoginScreen } from './login';
import { GameScreen } from './game'

export const ShengJiApp = () => {
    const [isSetName, setIsSetName] = useState(false)    
    const [playerName, setPlayerName] = useState("")
    return <div>
        {
            isSetName ? <GameScreen playerName={playerName}></GameScreen> : <LoginScreen playerName={playerName} setPlayerName={setPlayerName} setIsSetName={setIsSetName}></LoginScreen>
        }
    </div>
}
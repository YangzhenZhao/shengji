import { useEffect } from "react";
import Phaser from "phaser"
import { GameScene } from './game_scene';

interface GameScreenProps {
    playerName: string
}

export const GameScreen = ({playerName}: GameScreenProps) => {    
    useEffect(() => {
        const screenWidth = document.documentElement.clientWidth;
        const screenHeight =  document.documentElement.clientHeight;
        const gameScence = new GameScene(playerName)
        const config = {
            type: Phaser.AUTO,
            width: screenWidth,
            height: screenHeight,
            parent: 'phaser-example',
            scene: [gameScence]
        };
        new Phaser.Game(config);
    }, [])
    
    return (
        <div id="phaser-example"></div>
    )
}
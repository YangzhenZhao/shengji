import { GameScene } from './game_scene';

export class Match {
    public myTeamRound: string
    public oppositeTeamRound: string
    private gameScene: GameScene
    constructor(gameScene: GameScene) {
        this.myTeamRound = "2"
        this.oppositeTeamRound = "2"
        this.gameScene = gameScene
    }
}

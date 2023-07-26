export interface UserInfo {
    name: string;
    admin: boolean;
    expireTime: number;
}

export interface LobbyInfo {
    id: string;
    host: string;
    players: string[];
    started: boolean;
}

export interface Card {
    color: number;
    value: number;
}

export interface BidInfo {
    id: string;
    done: boolean;
    players: string[];
    hand: Card[];
    passed: boolean[];
    currentBidder: number;
    bid: number;
}

export interface GameInfo {
    id: string;
    done: boolean;
    players: string[];
    hand: Card[];
    table: Card[];
    currentPlayer: number;
    trump: number;
    bid: number;
    bidWinner: number;
}

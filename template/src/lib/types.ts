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
    hands: Card[][];
    widow: Card[];
    passed: boolean[];
    currentBidder: number;
    bid: number;
}

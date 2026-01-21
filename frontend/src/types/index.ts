// Types matching backend models

export interface RankedUser {
    rank: number;
    id: number;
    username: string;
    rating: number;
}

export interface LeaderboardResponse {
    users: RankedUser[];
    total: number;
    limit: number;
    offset: number;
}

export interface SearchResponse {
    users: RankedUser[];
    count: number;
}

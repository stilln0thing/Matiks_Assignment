import client from './client';
import { LeaderboardResponse, SearchResponse, RankedUser } from '../types';

export const getLeaderboard = async (
    limit: number = 50,
    offset: number = 0
): Promise<LeaderboardResponse> => {
    const response = await client.get<LeaderboardResponse>('/leaderboard', {
        params: { limit, offset },
    });
    return response.data;
};

export const searchUsers = async (query: string): Promise<SearchResponse> => {
    const response = await client.get<SearchResponse>('/search', {
        params: { q: query },
    });
    return response.data;
};

export const getUserRank = async (userId: number): Promise<RankedUser> => {
    const response = await client.get<RankedUser>(`/user/${userId}/rank`);
    return response.data;
};

export const updateRating = async (
    userId: number,
    rating: number
): Promise<void> => {
    await client.post('/rating', { user_id: userId, rating });
};

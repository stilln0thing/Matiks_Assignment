import React, { useState, useEffect, useCallback } from 'react';
import {
    View,
    FlatList,
    Text,
    StyleSheet,
    ActivityIndicator,
    RefreshControl,
} from 'react-native';
import { SafeAreaView } from 'react-native-safe-area-context';
import { getLeaderboard } from '../api/leaderboard';
import { RankedUser } from '../types';
import { LeaderboardItem } from '../components/LeaderboardItem';
import { colors, spacing, typography } from '../theme';

const LIMIT = 50;

export const LeaderboardScreen: React.FC = () => {
    const [users, setUsers] = useState<RankedUser[]>([]);
    const [total, setTotal] = useState(0);
    const [loading, setLoading] = useState(true);
    const [refreshing, setRefreshing] = useState(false);
    const [loadingMore, setLoadingMore] = useState(false);
    const [error, setError] = useState<string | null>(null);

    const fetchLeaderboard = async (offset: number = 0, isRefresh: boolean = false) => {
        try {
            if (isRefresh) {
                setRefreshing(true);
            } else if (offset === 0) {
                setLoading(true);
            } else {
                setLoadingMore(true);
            }
            setError(null);

            const response = await getLeaderboard(LIMIT, offset);

            if (offset === 0 || isRefresh) {
                setUsers(response.users);
            } else {
                setUsers(prev => [...prev, ...response.users]);
            }
            setTotal(response.total);
        } catch (err) {
            setError('Failed to load leaderboard. Please try again.');
            console.error('Leaderboard error:', err);
        } finally {
            setLoading(false);
            setRefreshing(false);
            setLoadingMore(false);
        }
    };

    useEffect(() => {
        fetchLeaderboard();
    }, []);

    const handleRefresh = useCallback(() => {
        fetchLeaderboard(0, true);
    }, []);

    const handleLoadMore = useCallback(() => {
        if (loadingMore || users.length >= total) return;
        fetchLeaderboard(users.length);
    }, [loadingMore, users.length, total]);

    const renderItem = useCallback(({ item }: { item: RankedUser }) => (
        <LeaderboardItem user={item} />
    ), []);

    const renderHeader = () => (
        <View style={styles.header}>
            <Text style={styles.title}>🏆 Leaderboard</Text>
            <Text style={styles.subtitle}>
                {total.toLocaleString()} players competing
            </Text>
        </View>
    );

    const renderFooter = () => {
        if (!loadingMore) return null;
        return (
            <View style={styles.footer}>
                <ActivityIndicator size="small" color={colors.primary} />
            </View>
        );
    };

    const renderEmpty = () => (
        <View style={styles.emptyContainer}>
            <Text style={styles.emptyText}>No players yet</Text>
        </View>
    );

    if (loading) {
        return (
            <SafeAreaView style={styles.container}>
                <View style={styles.loadingContainer}>
                    <ActivityIndicator size="large" color={colors.primary} />
                    <Text style={styles.loadingText}>Loading leaderboard...</Text>
                </View>
            </SafeAreaView>
        );
    }

    if (error) {
        return (
            <SafeAreaView style={styles.container}>
                <View style={styles.errorContainer}>
                    <Text style={styles.errorText}>{error}</Text>
                </View>
            </SafeAreaView>
        );
    }

    return (
        <SafeAreaView style={styles.container} edges={['top']}>
            <FlatList
                data={users}
                renderItem={renderItem}
                keyExtractor={(item) => item.id.toString()}
                ListHeaderComponent={renderHeader}
                ListFooterComponent={renderFooter}
                ListEmptyComponent={renderEmpty}
                onEndReached={handleLoadMore}
                onEndReachedThreshold={0.3}
                refreshControl={
                    <RefreshControl
                        refreshing={refreshing}
                        onRefresh={handleRefresh}
                        tintColor={colors.primary}
                        colors={[colors.primary]}
                    />
                }
                showsVerticalScrollIndicator={false}
                contentContainerStyle={styles.listContent}
            />
        </SafeAreaView>
    );
};

const styles = StyleSheet.create({
    container: {
        flex: 1,
        backgroundColor: colors.background,
    },
    listContent: {
        paddingBottom: spacing.xl,
    },
    header: {
        padding: spacing.lg,
        paddingBottom: spacing.md,
    },
    title: {
        color: colors.text,
        fontSize: typography.h1.fontSize,
        fontWeight: typography.h1.fontWeight,
    },
    subtitle: {
        color: colors.textSecondary,
        fontSize: typography.caption.fontSize,
        marginTop: spacing.xs,
    },
    loadingContainer: {
        flex: 1,
        justifyContent: 'center',
        alignItems: 'center',
    },
    loadingText: {
        color: colors.textSecondary,
        marginTop: spacing.md,
        fontSize: typography.body.fontSize,
    },
    errorContainer: {
        flex: 1,
        justifyContent: 'center',
        alignItems: 'center',
        padding: spacing.lg,
    },
    errorText: {
        color: colors.error,
        fontSize: typography.body.fontSize,
        textAlign: 'center',
    },
    emptyContainer: {
        flex: 1,
        justifyContent: 'center',
        alignItems: 'center',
        padding: spacing.xl,
    },
    emptyText: {
        color: colors.textSecondary,
        fontSize: typography.body.fontSize,
    },
    footer: {
        paddingVertical: spacing.md,
    },
});

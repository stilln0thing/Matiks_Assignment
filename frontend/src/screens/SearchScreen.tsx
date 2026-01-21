import React, { useState, useCallback, useEffect } from 'react';
import {
    View,
    FlatList,
    Text,
    StyleSheet,
    ActivityIndicator,
} from 'react-native';
import { SafeAreaView } from 'react-native-safe-area-context';
import { searchUsers } from '../api/leaderboard';
import { RankedUser } from '../types';
import { SearchBar } from '../components/SearchBar';
import { LeaderboardItem } from '../components/LeaderboardItem';
import { colors, spacing, typography } from '../theme';

const DEBOUNCE_MS = 300;

export const SearchScreen: React.FC = () => {
    const [query, setQuery] = useState('');
    const [results, setResults] = useState<RankedUser[]>([]);
    const [loading, setLoading] = useState(false);
    const [error, setError] = useState<string | null>(null);
    const [hasSearched, setHasSearched] = useState(false);

    const performSearch = useCallback(async (searchQuery: string) => {
        if (!searchQuery.trim()) {
            setResults([]);
            setHasSearched(false);
            return;
        }

        try {
            setLoading(true);
            setError(null);
            const response = await searchUsers(searchQuery);
            setResults(response.users);
            setHasSearched(true);
        } catch (err) {
            setError('Search failed. Please try again.');
            console.error('Search error:', err);
        } finally {
            setLoading(false);
        }
    }, []);

    // Debounced search
    useEffect(() => {
        const timer = setTimeout(() => {
            performSearch(query);
        }, DEBOUNCE_MS);

        return () => clearTimeout(timer);
    }, [query, performSearch]);

    const renderItem = useCallback(({ item }: { item: RankedUser }) => (
        <LeaderboardItem user={item} />
    ), []);

    const renderEmpty = () => {
        if (loading) {
            return (
                <View style={styles.centerContainer}>
                    <ActivityIndicator size="large" color={colors.primary} />
                    <Text style={styles.statusText}>Searching...</Text>
                </View>
            );
        }

        if (error) {
            return (
                <View style={styles.centerContainer}>
                    <Text style={styles.errorText}>{error}</Text>
                </View>
            );
        }

        if (!hasSearched) {
            return (
                <View style={styles.centerContainer}>
                    <Text style={styles.hintEmoji}>🔍</Text>
                    <Text style={styles.hintText}>
                        Search for players by username
                    </Text>
                    <Text style={styles.hintSubtext}>
                        Results will appear as you type
                    </Text>
                </View>
            );
        }

        return (
            <View style={styles.centerContainer}>
                <Text style={styles.hintEmoji}>🤷</Text>
                <Text style={styles.noResultsText}>No players found</Text>
                <Text style={styles.hintSubtext}>Try a different search term</Text>
            </View>
        );
    };

    return (
        <SafeAreaView style={styles.container} edges={['top']}>
            <View style={styles.header}>
                <Text style={styles.title}>🔎 Search</Text>
                <Text style={styles.subtitle}>Find players by username</Text>
            </View>

            <SearchBar value={query} onChangeText={setQuery} />

            {results.length > 0 && (
                <Text style={styles.resultsCount}>
                    {results.length} result{results.length !== 1 ? 's' : ''} found
                </Text>
            )}

            <FlatList
                data={results}
                renderItem={renderItem}
                keyExtractor={(item) => item.id.toString()}
                ListEmptyComponent={renderEmpty}
                showsVerticalScrollIndicator={false}
                contentContainerStyle={[
                    styles.listContent,
                    results.length === 0 && styles.listEmpty,
                ]}
                keyboardShouldPersistTaps="handled"
            />
        </SafeAreaView>
    );
};

const styles = StyleSheet.create({
    container: {
        flex: 1,
        backgroundColor: colors.background,
    },
    header: {
        padding: spacing.lg,
        paddingBottom: spacing.sm,
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
    resultsCount: {
        color: colors.textSecondary,
        fontSize: typography.caption.fontSize,
        marginHorizontal: spacing.lg,
        marginBottom: spacing.sm,
    },
    listContent: {
        paddingBottom: spacing.xl,
    },
    listEmpty: {
        flex: 1,
    },
    centerContainer: {
        flex: 1,
        justifyContent: 'center',
        alignItems: 'center',
        padding: spacing.xl,
    },
    hintEmoji: {
        fontSize: 48,
        marginBottom: spacing.md,
    },
    hintText: {
        color: colors.textSecondary,
        fontSize: typography.body.fontSize,
        textAlign: 'center',
    },
    hintSubtext: {
        color: colors.textMuted,
        fontSize: typography.caption.fontSize,
        textAlign: 'center',
        marginTop: spacing.xs,
    },
    statusText: {
        color: colors.textSecondary,
        marginTop: spacing.md,
        fontSize: typography.body.fontSize,
    },
    errorText: {
        color: colors.error,
        fontSize: typography.body.fontSize,
        textAlign: 'center',
    },
    noResultsText: {
        color: colors.textSecondary,
        fontSize: typography.body.fontSize,
        textAlign: 'center',
    },
});

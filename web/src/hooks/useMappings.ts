import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query';
import { addMapping, deleteMapping, fetchMappings, PathUrlMapping, updateMapping } from '../api/mappings';

export const useMappings = (offset: number, limit: number) => {
    return useQuery({
        queryKey: ['mappings', offset, limit],
        queryFn: async () => { const rlt = await fetchMappings(offset, limit); console.log(rlt); return rlt; },
        staleTime: 5 * 60 * 1000, // 5 minutes
        refetchInterval: 30 * 1000, // 30 seconds
        retry: 2,
        retryDelay: 1000,
    });
};

export const useAddMapping = () => {
    const queryClient = useQueryClient();
    return useMutation({
        mutationFn: addMapping,
        onSuccess: () => {
            queryClient.invalidateQueries({ queryKey: ['mappings'] });
        },
    });
};

export const useUpdateMapping = () => {
    const queryClient = useQueryClient();
    return useMutation({
        mutationFn: (mapping: Omit<PathUrlMapping, 'mapper' | 'usecount'>) => updateMapping(mapping),
        onSuccess: () => {
            queryClient.invalidateQueries({ queryKey: ['mappings'] });
        },
    });
};

export const useDeleteMapping = () => {
    const queryClient = useQueryClient();
    return useMutation({
        mutationFn: deleteMapping,
        onSuccess: () => {
            queryClient.invalidateQueries({ queryKey: ['mappings'] });
        },
    });
};

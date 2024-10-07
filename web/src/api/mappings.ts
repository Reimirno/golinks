const API_URL = '/api/go/';

export interface PathUrlMapping {
    path: string;
    url: string;
    mapper: string;
    usecount: number;
}

export const fetchMappings = async (offset = 0, limit = 10) => {
    const response = await fetch(`${API_URL}?offset=${offset}&limit=${limit}`);
    if (!response.ok) {
        throw new Error('Network response was not ok');
    }
    return response.json();
};

export const addMapping = async (mapping: Omit<PathUrlMapping, 'mapper' | 'usecount'>) => {
    const response = await fetch(API_URL, {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json'
        },
        body: JSON.stringify(mapping)
    });
    if (!response.ok) {
        throw new Error('Failed to add mapping');
    }
    return response.json();
};

export const updateMapping = async (mapping: Omit<PathUrlMapping, 'mapper' | 'usecount'>) => {
    const response = await fetch(API_URL, {
        method: 'PUT',
        headers: {
            'Content-Type': 'application/json'
        },
        body: JSON.stringify(mapping)
    });
    if (!response.ok) {
        throw new Error('Failed to update mapping');
    }
    return response.json();
};

export const deleteMapping = async (path: string) => {
    const response = await fetch(`${API_URL}/${path}`, {
        method: 'DELETE'
    });
    if (!response.ok) {
        throw new Error('Failed to delete mapping');
    }
};

import React, { useState } from 'react';
import { useAddMapping } from '../hooks/useMappings';

const MappingForm: React.FC = () => {
    const addMutation = useAddMapping();
    const [formData, setFormData] = useState({
        path: '',
        url: '',
        mapper: '',
    });

    const handleChange = (e: React.ChangeEvent<HTMLInputElement>) => {
        setFormData(prev => ({ ...prev, [e.target.name]: e.target.value }));
    };

    const handleSubmit = (e: React.FormEvent) => {
        e.preventDefault();
        addMutation.mutate(formData);
        setFormData({ path: '', url: '', mapper: '' });
    };

    return (
        <form onSubmit={handleSubmit} className="p-4 space-y-4">
            <div className="form-control">
                <label className="label">Path</label>
                <input
                    className="input input-bordered"
                    name="path"
                    value={formData.path}
                    onChange={handleChange}
                    required
                />
            </div>
            <div className="form-control">
                <label className="label">URL</label>
                <input
                    className="input input-bordered"
                    name="url"
                    value={formData.url}
                    onChange={handleChange}
                    required
                />
            </div>
            <div className="form-control">
                <label className="label">Mapper</label>
                <input
                    className="input input-bordered"
                    name="mapper"
                    value={formData.mapper}
                    onChange={handleChange}
                    required
                />
            </div>
            <button className="btn btn-primary" type="submit" disabled={addMutation.isPending}>
                Add Mapping
            </button>
        </form>
    );
};

export default MappingForm;

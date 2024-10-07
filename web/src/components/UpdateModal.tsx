import React, { useState } from 'react';
import { PathUrlMapping } from '../api/mappings';
import { useUpdateMapping } from '../hooks/useMappings';

interface UpdateModalProps {
    mapping: PathUrlMapping;
    onClose: () => void;
}

const UpdateModal: React.FC<UpdateModalProps> = ({ mapping, onClose }) => {
    const updateMutation = useUpdateMapping();
    const [formData, setFormData] = useState({
        path: mapping.path,
        url: mapping.url,
        mapper: mapping.mapper,
    });

    const handleChange = (e: React.ChangeEvent<HTMLInputElement>) => {
        setFormData(prev => ({ ...prev, [e.target.name]: e.target.value }));
    };

    const handleSubmit = (e: React.FormEvent) => {
        e.preventDefault();
        updateMutation.mutate(formData);
        onClose();
    };

    return (
        <div className="modal modal-open">
            <div className="modal-box">
                <h3 className="font-bold text-lg">Update Mapping</h3>
                <form onSubmit={handleSubmit} className="space-y-4">
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
                    <div className="modal-action">
                        <button className="btn" type="submit" disabled={updateMutation.isPending}>
                            Save
                        </button>
                        <button className="btn btn-outline" onClick={onClose} type="button">
                            Cancel
                        </button>
                    </div>
                </form>
            </div>
        </div>
    );
};

export default UpdateModal;

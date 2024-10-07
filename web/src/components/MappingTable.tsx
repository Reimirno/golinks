import 'ag-grid-community/styles/ag-grid.css';
import 'ag-grid-community/styles/ag-theme-alpine.css';
import { AgGridReact } from 'ag-grid-react';
import React, { useMemo, useState } from 'react';
import { useDeleteMapping, useMappings } from '../hooks/useMappings';
import UpdateModal from './UpdateModal.tsx';

const MappingTable: React.FC = () => {
    const [offset, setOffset] = useState(0);
    const limit = 10;
    const { data, isLoading, isError, error } = useMappings(offset, limit);
    const deleteMutation = useDeleteMapping();
    const [selectedMapping, setSelectedMapping] = useState(null);

    const columns = useMemo(() => [
        { headerName: 'Path', field: 'path', sortable: true },
        { headerName: 'URL', field: 'url', sortable: true },
        { headerName: 'Mapper', field: 'mapper', sortable: true },
        { headerName: 'Use Count', field: 'usecount', sortable: true },
        {
            headerName: 'Actions',
            field: 'id',
            cellRendererFramework: (params: any) => (
                <div className="flex space-x-2">
                    <button
                        className="btn btn-xs"
                        onClick={() => setSelectedMapping(params.data)}
                    >
                        Edit
                    </button>
                    <button
                        className="btn btn-xs btn-error"
                        onClick={() => deleteMutation.mutate(params.data.id)}
                    >
                        Delete
                    </button>
                </div>
            ),
        },
    ], [deleteMutation]);

    return (
        <div className="ag-theme-alpine" style={{ height: 500, width: '100%' }}>
            {isLoading ? (
                <div>Loading...</div>
            ) : isError ? (
                <div>Error: {error.message}</div>
            ) : (
                <>
                    <AgGridReact
                        rowData={data}
                        columnDefs={columns}
                        pagination={true}
                        paginationPageSize={limit}
                    />
                    {/* <div className="flex justify-between mt-2">
                        <button
                            className="btn btn-sm"
                            onClick={() => setOffset(Math.max(offset - limit, 0))}
                            disabled={offset === 0}
                        >
                            Previous
                        </button>
                        <button
                            className="btn btn-sm"
                            onClick={() => setOffset(offset + limit)}
                            disabled={data?.mappings.length < limit}
                        >
                            Next
                        </button>
                    </div> */}
                    {selectedMapping && (
                        <UpdateModal
                            mapping={selectedMapping}
                            onClose={() => setSelectedMapping(null)}
                        />
                    )}
                </>
            )}
        </div>
    );
};

export default MappingTable;

<script lang="ts">
	import type { ColDef, GridOptions } from 'ag-grid-community';
	import 'ag-grid-community/styles/ag-grid.css';
	import 'ag-grid-community/styles/ag-theme-alpine.css';
	import AgGridSvelte from 'ag-grid-svelte';
	import { onMount } from 'svelte';

	import ActionsCellRenderer from './ActionsCellRenderer.svelte';
	import MappingModal from './MappingModal.svelte';

	interface Mapping {
		path: string;
		url: string;
		mapper: string;
		usecount: number;
	}

	let gridOptions: GridOptions = {
		context: {
			deleteMapping,
			openEditModal
		}
	};

	let columnDefs: ColDef[] = [
		{ field: 'path', headerName: 'Path', sortable: true, filter: true },
		{ field: 'url', headerName: 'URL', sortable: true, filter: true },
		{ field: 'mapper', headerName: 'Mapper', sortable: true, filter: true },
		{ field: 'usecount', headerName: 'Use Count', sortable: true, filter: true },
		{
			headerName: 'Actions',
			cellRenderer: ActionsCellRenderer,
			editable: false,
			filter: false,
			sortable: false
		}
	];

	let rowData: Mapping[] = [];
	let currentPage = 1;
	let pageSize = 10;
	let totalPages = 1;
	let sortModel = [];

	let isModalOpen = false;
	let modalMapping: Mapping = { path: '', url: '', mapper: '', usecount: 0 };
	let isEdit = false;

	async function fetchMappings(page = 1) {
		let offset = (page - 1) * pageSize;
		let limit = pageSize;
		let url = `/go?offset=${offset}&limit=${limit}`;

		// if (sortModel.length > 0) {
		// 	const sortField = sortModel[0].colId;
		// 	const sortOrder = sortModel[0].sort;
		// 	url += `&sortField=${sortField}&sortOrder=${sortOrder}`;
		// }

		const response = await fetch(url);
		const data = await response.json();

		rowData = data.mappings;
		const totalItems = data.totalItems;
		totalPages = Math.ceil(totalItems / pageSize);
	}

	onMount(() => {
		fetchMappings(currentPage);
	});

	function deleteMapping(path: string) {
		fetch(`/api/mappings/${encodeURIComponent(path)}`, {
			method: 'DELETE'
		}).then((response) => {
			if (response.ok) {
				fetchMappings(currentPage);
			} else {
				alert('Failed to delete mapping.');
			}
		});
	}

	function openAddModal() {
		isModalOpen = true;
		isEdit = false;
		modalMapping = { path: '', url: '', mapper: '', usecount: 0 };
	}

	function openEditModal(mapping) {
		isModalOpen = true;
		isEdit = true;
		modalMapping = { ...mapping };
	}

	function handleModalSave(event) {
		const mapping = event.detail;
		const method = isEdit ? 'PUT' : 'POST';
		const url = isEdit ? `/api/mappings/${encodeURIComponent(mapping.path)}` : '/api/mappings';

		fetch(url, {
			method,
			headers: { 'Content-Type': 'application/json' },
			body: JSON.stringify(mapping)
		}).then((response) => {
			if (response.ok) {
				fetchMappings(currentPage);
			} else {
				alert('Failed to save mapping.');
			}
		});
	}

	function nextPage() {
		if (currentPage < totalPages) {
			currentPage++;
			fetchMappings(currentPage);
		}
	}

	function prevPage() {
		if (currentPage > 1) {
			currentPage--;
			fetchMappings(currentPage);
		}
	}

	function onSortChanged(event) {
		sortModel = event.api.getSortModel();
		fetchMappings(currentPage);
	}
</script>

<button class="btn btn-primary mb-4" on:click={openAddModal}>Add Mapping</button>

<MappingModal
	bind:isOpen={isModalOpen}
	bind:mapping={modalMapping}
	{isEdit}
	on:save={handleModalSave}
/>

<AgGridSvelte
	class="ag-theme-alpine"
	{rowData}
	{columnDefs}
	{gridOptions}
	{onSortChanged}
	style="width: 100%; height: 500px;"
/>

<div class="mt-4 flex justify-between">
	<button class="btn" on:click={prevPage} disabled={currentPage <= 1}> Previous </button>
	<span>
		Page {currentPage} of {totalPages}
	</span>
	<button class="btn" on:click={nextPage} disabled={currentPage >= totalPages}> Next </button>
</div>

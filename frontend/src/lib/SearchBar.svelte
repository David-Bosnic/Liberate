<script lang="ts">
    import { liberateState } from "./liberateStore.svelte";
    import LargeSearchBar from "./LargeSearchBar.svelte";
    import CompactSearchBar from "./CompactSearchBar.svelte";
    import ResultsList from "./ResultsList.svelte";

    let hasResults = $derived(liberateState.links.length > 0);
</script>

{#if !hasResults}
    <LargeSearchBar />
{:else}
    <CompactSearchBar />
{/if}

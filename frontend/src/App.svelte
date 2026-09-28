<script lang="ts">
    import SearchBar from "./lib/SearchBar.svelte";
    import ResultsList from "./lib/ResultsList.svelte";
</script>

<div class="flex flex-col justify-center items-center">
    <SearchBar />
    <div class="w-full max-w-md">
        <ResultsList />
    </div>
</div>

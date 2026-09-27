<script lang="ts">
    import { liberateState } from "./liberateStore.svelte";

    let hasResults = $derived(liberateState.links.length > 0);
</script>

{#if hasResults}
    {#each liberateState.links as link}
        <a class="" href={link} target="_blank" rel="noopener noreferrer">
            <div class="p-2 gap-2 text-secondary-300 rounded-2xl border-2">
                {link}
            </div>
        </a>
    {/each}
{/if}

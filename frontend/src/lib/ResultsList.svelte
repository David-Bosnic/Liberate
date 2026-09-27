<script lang="ts">
    import { liberateState } from "./liberateStore.svelte";

    let hasResults = $derived(liberateState.links.length > 0);
</script>

{#if hasResults}
    {#each liberateState.links as link}
        <a
            class="block text-left mb-4"
            href={link}
            target="_blank"
            rel="noopener noreferrer"
        >
            <div class="text-secondary-300">
                {link}
            </div>
        </a>
    {/each}
{/if}

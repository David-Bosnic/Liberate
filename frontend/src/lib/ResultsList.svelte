<script lang="ts">
    import { liberateState } from "./liberateStore.svelte";

    let hasResults = $derived(liberateState.links.length > 0);
    
    function getRatingColor(rating: string): string {
        const ratingMap: Record<string, string> = {
            "1": "text-red-500",
            "2": "text-orange-500",
            "3": "text-yellow-500",
            "4": "text-green-500",
            "5": "text-emerald-500",
        };
        return ratingMap[rating] || "text-gray-500";
    }
    
    function getRatingLabel(rating: string): string {
        const labelMap: Record<string, string> = {
            "1": "Poor",
            "2": "Bad",
            "3": "OK",
            "4": "Good",
            "5": "Great",
        };
        return labelMap[rating] || rating;
    }
</script>

{#if hasResults}
    {#each liberateState.links as link}
        <a
            class="block text-left mb-4"
            href={link.url}
            target="_blank"
            rel="noopener noreferrer"
        >
            <div class="flex justify-between items-center">
                <div class="text-secondary-300">
                    {link.url}
                </div>
                <div class={`text-sm font-medium ${getRatingColor(link.rating)}`}>
                    Rating: {getRatingLabel(link.rating)} ({link.rating})
                </div>
            </div>
        </a>
    {/each}
{/if}

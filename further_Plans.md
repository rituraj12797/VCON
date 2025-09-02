# VCON System Architecture: Content Addressable Store (CAS) and Hydration

This document outlines the core storage architecture for the VCON project, focusing on how versioned documents are stored efficiently and retrieved quickly.

## 1. Core Architecture: Two-Collection Model

Our architecture is based on two distinct but connected MongoDB collections: a `documents` collection for metadata and a `content_strings` collection that acts as a global Content Addressable Store (CAS).

### a. `documents` Collection

-   **Purpose:** Stores the metadata and version tree structure for each file managed by VCON. Each document in this collection represents a single versioned file.
-   **Schema:** Contains fields like `_id` (a standard `ObjectID`), `title`, and `NodeArray`.
-   **Key Feature:** The `Node` structs within the `NodeArray` do **not** store the actual text content. Instead, they store arrays of integers (`[]int`) in their `FileContent` and `DeltaInstructions` fields. These integers are pointers to the content in the CAS.

### b. `content_strings` Collection (The CAS)

-   **Purpose:** Acts as a single, global "dictionary" for the entire application. It ensures that every unique line of text across all files and all versions is stored only once.
-   **Schema:** Each document in this collection is a simple mapping.
    -   `_id` (`int`): A unique integer identifier for the string. This is the primary key.
    -   `content` (`string`): The actual line of text.
-   **Indexing:** A **unique index** will be created on the `content` field. This guarantees that no string can be inserted more than once and makes lookups by content (for "interning" new strings) extremely fast.

## 2. The Hydration Strategy: Efficiently Loading Documents

When a user wants to open or work with a file, we don't want to make database calls every time they switch versions. Instead, we will "hydrate" the document once upon loading.

This is a multi-step process designed for maximum performance:

**Step 1: Fetch the "Dry" Document**
-   We perform a single query to the `documents` collection to retrieve the specific document object by its `_id` or `title`. This object contains the full version tree but only the integer IDs for content.

**Step 2: Collect All Unique Content IDs**
-   Once the document is in the application's memory, we perform a fast, in-memory loop through its entire `NodeArray`.
-   We gather every integer ID from every `FileContent` and `DeltaInstructions` array.
-   These IDs are collected into a `set` (in Go, a `map[int]struct{}`) to create a list of all unique content strings required to render *any* version of this document.

**Step 3: Bulk Fetch from CAS**
-   We perform **one single, highly efficient bulk query** against the `content_strings` collection.
-   This query uses MongoDB's `$in` operator with the unique set of IDs collected in Step 2.
-   Example: `db.content_strings.find({ _id: { $in: [101, 45, 800, 2, ...] } })`

**Step 4: Build In-Memory Cache**
-   The results from the bulk query are used to populate a simple in-memory `map[int]string`.
-   This map serves as a local, temporary, and extremely fast "global store" specifically for the loaded document.

## 3. Rendering and Performance

-   **Zero DB Calls for Rendering:** Once the document is hydrated and the in-memory cache is built, rendering any version (by applying deltas or displaying snapshots) is blazing fast. All content lookups are simple map lookups in memory, requiring zero further database interaction.
-   **Efficiency:** This strategy minimizes database load by replacing potentially thousands of small, slow queries with one initial, indexed bulk query.

## 4. Benefits of this Architecture

-   **Massive Deduplication:** Saves enormous amounts of storage by never storing the same line of text twice.
-   **High Performance:** Leverages database indexing for fast lookups and minimizes network round-trips through bulk operations.
-   **Scalability:** The CAS can scale to billions of unique strings without hitting document size limits.
-   **Data Integrity:** The `NodeArray` acts as the single source of truth for a document's structure, preventing the data corruption risks associated with storing redundant
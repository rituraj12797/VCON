# VCON System Architecture: Hash-Based CAS and Hydration Strategy

This document outlines the core storage architecture for the VCON project. It is based on a highly scalable, hash-based Content Addressable Store (CAS) model, similar to systems like Git.

## 1. Core Architecture: Two-Collection Model

Our architecture is based on two distinct but connected MongoDB collections: a `documents` collection for metadata and a `content_store` collection that acts as a global CAS.

### a. `documents` Collection

-   **Purpose:** Stores the metadata and version tree structure for each file managed by VCON.
-   **Schema:** Contains fields like `_id` (a standard `ObjectID`), `title`, and `NodeArray`.
-   **Key Feature:** The `Node` structs within the `NodeArray` do **not** store the actual text content. Instead, they store arrays of **strings** (`[]string`) in their `FileContent` and `DeltaInstructions` fields. These strings are the SHA-256 hashes of the content.

### b. `content_store` Collection (The Global CAS)

-   **Purpose:** Acts as a single, global "dictionary" for the entire application. It ensures that every unique line of text across all files and all versions is stored only once.
-   **Schema:** Each document in this collection is a simple mapping from a hash to its content.
    -   `_id` (`string`): The **SHA-256 hash** of the content, stored as a hex string. This is the primary key.
    -   `content` (`string`): The actual line of text.
-   **Benefit:** This design makes "interning" new content stateless. To get the ID for a new line of text, the application simply calculates its SHA-256 or MURMUR hash. No database query is needed to generate an ID.

## 2. The Hydration and Caching Strategy

When a user opens a file, we perform a one-time "hydration" process to load all necessary content into a high-speed cache for rendering.

**Step 1: Fetch the "Dry" Document**
-   Perform a single query to the `documents` collection to retrieve the specific document object. This object contains the full version tree but only the SHA-256 hashes for content.

**Step 2: Collect All Unique Hashes**
-   Once the document is in memory, perform a fast loop through its `NodeArray` to gather every unique hash from all `FileContent` and `DeltaInstructions` arrays into a `set`.

**Step 3: Check the Redis Cache**
-   Before querying the main database, we first check our Redis cache for the required hashes. This reduces the load on the main database for frequently accessed content.

**Step 4: Bulk Fetch from Main Database (Cache Misses)**
-   For any hashes that were not found in Redis (a "cache miss"), we perform **one single, highly efficient bulk query** against the `content_store` collection in MongoDB.
-   This query uses the `$in` operator with the list of missing hashes.
-   Example: `db.content_store.find({ _id: { $in: ["a1b2...", "c3d4...", ...] } })`

**Step 5: Populate Cache and Build In-Memory Map**
-   The results from the database are used to populate the Redis cache for future requests.
-   All the content (from both Redis and the DB) is then loaded into a simple in-memory `map[string]string` within the application for the fastest possible access during the user's session.

## 3. Benefits of this Architecture

-   **Stateless ID Generation:** Massively simplifies client-side logic and scales better in distributed environments.
-   **Guaranteed Global Deduplication:** Identical content is guaranteed to have the same ID and is stored only once.
-   **Built-in Data Integrity:** The hash acts as a fingerprint, ensuring the content has not been corrupted.
-   **High Performance:** The multi-layered caching (Redis + in-memory map) and bulk database operations ensure that rendering versions is extremely fast after the initial load.



Remaining now 
0.  (a) Global storages for mapping title vs document to keep dcuments into RAM as cached
    (b) global primary memory level map for hash vs strings of all the subsets of CAS of diffrent documents 

    


1. API ( Get document X ) - this will load the dcument from Db ans tore in a map defined in main.go and store the contents of it's subset of CAS into our in memory global CAS

2. API ( Create Document ) - workflow in copy 

3. APi ( Add version x) - workflow in copy 

4. API ( Get version X ) - workflow in copy

5. Multiprocesing Hasher 

6. Codin the main.go the ( REST ) Server 

7. Refine the already existing get version x to work and return an array f hash from a destined version to a give version working on our defined node structure now 


package Scrunching;

import com.mongodb.MongoClient;
import com.mongodb.client.MongoCollection;
import com.mongodb.client.MongoDatabase;
import org.bson.Document;
import com.mongodb.client.FindIterable;
import com.mongodb.Block;


import java.io.*;
import java.nio.file.Path;
import java.nio.file.Paths;
import java.util.*;
import java.util.List;

/**
 * Created by Joel Shapiro on 10/27/15.
 * Part of the CodeCollaborate project
 * <p/>
 * THIS ASSUMES FILES WILL FIT IN MEMORY
 */

public class Scrunching {

    private static DiffMatchPatch differ;

    public static void main(String[] args) {
        if (args.length < 1) {
            throw new RuntimeException("No fileId supplied");
        }
        String fileId = args[0];

        differ = new DiffMatchPatch();
        differ.Diff_Timeout = 2.0f; // Ensure it won't fail from timing issue

        MongoClient mongoClient = new MongoClient();
        MongoDatabase db = mongoClient.getDatabase("CodeCollaborate");

        Document dbFileEntry = db.getCollection("Files").find(new Document("_id", fileId)).first();
        String rawPath = "files/" + dbFileEntry.get("project") + "/" + dbFileEntry.get("relative_path") + dbFileEntry.get("name");

        Path path = Paths.get(rawPath);
        String filePath = path.toAbsolutePath().toString();

        BufferedReader fileIn = null;

//        System.out.println(filePath);

        try {
            fileIn = new BufferedReader(new FileReader(filePath));
        } catch (FileNotFoundException e) {
            System.out.println("File not found");
            System.exit(1);
        }

        StringBuilder stringBuilder = new StringBuilder();
        String line;

        try {
            while ((line = fileIn.readLine()) != null) {
                stringBuilder.append(line).append("\n");
            }
        } catch (IOException e) {
            System.out.println("Unknown error while reading file");
            System.exit(1);
        } finally {
            try {
                if (fileIn != null) {
                    fileIn.close();
                }
            } catch (IOException e) {

            }
        }

        String fileString = stringBuilder.toString();


        LinkedHashMap<String, String> patches = getPatches(fileId, db);

        String key;
        String patch;

        ArrayList<String> idsToRemove = new ArrayList<String>();

        for (Map.Entry<String, String> entry : patches.entrySet()) {
            key = entry.getKey();
            patch = entry.getValue();

            Object[] applied = scrunch(patch, fileString);
            if (((boolean[]) applied[1])[0]) {
                fileString = (String) applied[0];
                idsToRemove.add(key);
            } else {
                System.out.println("Can't apply patch: " + key + " to file: " + filePath);
                System.exit(1);
            }
        }

        MongoCollection changes = db.getCollection("Changes");
        // have to do a bad way because collection.remove is no longer a thing
        for (String id : idsToRemove) {
            if (!changes.deleteOne(new Document("_id", id)).wasAcknowledged()) {
                System.out.println("ERROR: unable to delete file");
            }
        }

        FileWriter fileOut = null;
        try {
            fileOut = new FileWriter(filePath, false);
            fileOut.write(fileString);
        } catch (IOException e) {

        } finally {
            try {
                if (fileOut != null) {
                    fileOut.close();
                }
            } catch (IOException e) {

            }
        }
    }

    /**
     * @param fileId: id of the file to find
     * @return map of _ids to changes
     */
    private static LinkedHashMap<String, String> getPatches(String fileId, MongoDatabase db) {

        FindIterable<Document> iterable = db.getCollection("Changes").find(new Document("file", fileId)).sort(new Document("date", 1));

        final LinkedHashMap<String, String> patches = new LinkedHashMap<String, String>();

        iterable.forEach(new Block<Document>() {
            public void apply(final Document document) {
                patches.put((String) document.get("_id"), (String) document.get("changes"));
            }
        });

        return patches;
    }

    private static Object[] scrunch(String patch, String fileString) {
        List<DiffMatchPatch.Patch> patchList = null;

        try {
            patchList = differ.patch_fromText(patch);
        } catch (IllegalArgumentException e) {
            System.out.println("Unable to compile patch: \n" + patch);
            System.exit(1);
        }

        LinkedList<DiffMatchPatch.Patch> patchLinkedList = new LinkedList<DiffMatchPatch.Patch>(patchList);
        return differ.patch_apply(patchLinkedList, fileString);
    }

}

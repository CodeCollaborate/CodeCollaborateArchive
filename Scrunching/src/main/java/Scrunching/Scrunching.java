package Scrunching;

import com.mongodb.MongoClient;
import com.mongodb.client.MongoDatabase;
import org.bson.Document;
import com.mongodb.client.FindIterable;
import com.mongodb.Block;

import javax.sound.midi.SysexMessage;
import java.util.ArrayList;
import java.util.LinkedHashMap;
import java.util.Map;

/**
 * Created by Joel Shapiro on 10/27/15.
 * Part of the CodeCollaborate project
 */

public class Scrunching {

    private static diff_match_patch differ;

    public static void main(String[] args) {
        if (args.length < 1) {
            throw new RuntimeException("No fileId supplied");
        }
        String fileId = args[0];
        LinkedHashMap<String, String> patches = getPatches(fileId);

        differ = new diff_match_patch();
        differ.Diff_Timeout = 2.0f; // Ensure it won't fail from timing issue

//        for (Map.Entry<String, String> entry : patches.entrySet()){
//            System.out.println(entry.getKey() + ": " + entry.getValue());
//        }

        return;
    }

    /**
     *
     * @param fileId
     * @return map of _ids to changes
     */
    private static LinkedHashMap<String, String> getPatches(String fileId) {
        MongoClient mongoClient = new MongoClient();
        MongoDatabase db = mongoClient.getDatabase("CodeCollaborate");

        FindIterable<Document> iterable = db.getCollection("Changes").find(new Document("file", fileId)).sort(new Document("date", 1));

        final LinkedHashMap<String, String> patches = new LinkedHashMap<String, String>();

        iterable.forEach(new Block<Document>() {
            public void apply(final Document document) {
                patches.put((String) document.get("_id"), (String) document.get("changes"));
            }
        });

        return patches;
    }
}

package com.bullercodeworks.magopie;

import org.json.JSONException;
import org.json.JSONObject;

import java.util.ArrayList;
import java.util.HashMap;

import go.magopie.Magopie;

/**
 * Created by brbuller on 1/23/16.
 */
public class State {

  public String file = "/data/data/com.bullercodeworks.magopie/state";

  public String ServerURL = "";
  public String ApiToken = "";
  public ArrayList<Magopie.Torrent> results;
  public Magopie.Client client;
  public HashMap<String, String> sites;

  public State() {
    load();
    results = new ArrayList<>();
    sites = new HashMap<>();
    client = Magopie.NewClient(ServerURL, ApiToken);
    Magopie.SiteCollection s = UpdateSites();
    for(int i = 0; i < s.Length(); i++) {
      sites.put(s.Get(i).getID(), s.Get(i).getName());
    }
  }

  public Magopie.SiteCollection UpdateSites() {
    return client.ListSites();
  }

  public void save() {
    try {
      JSONObject jsonState = new JSONObject();
      jsonState.put("serverURL", ServerURL);
      jsonState.put("apiToken", ApiToken);
      Magopie.SaveToFile(file, jsonState.toString().getBytes());
    } catch(Exception e) { }
  }

  public void load() {
    String ret = "";
    byte[] res = Magopie.ReadFromFile(file);
    if(res != null && res.length > 0) {
      ret = new String(res);
      try {
        JSONObject jsonState = new JSONObject(ret);
        ServerURL = jsonState.getString("serverURL");
        ApiToken = jsonState.getString("apiToken");
      } catch(JSONException e) {}
    }
  }
}

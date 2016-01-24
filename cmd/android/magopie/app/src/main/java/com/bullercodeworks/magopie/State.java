package com.bullercodeworks.magopie;

import org.json.JSONException;
import org.json.JSONObject;

import java.util.ArrayList;

import go.magopie.Magopie;

/**
 * Created by brbuller on 1/23/16.
 */
public class State {

  public String file = "/data/data/com.bullercodeworks.magopie/state";

  public String ServerURL = "";
  public ArrayList<Magopie.Torrent> results;

  public State() {
    results = new ArrayList<>();
  }
  public void save() {
    try {
      JSONObject jsonState = new JSONObject();
      jsonState.put("serverURL", ServerURL);
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
      } catch(JSONException e) {}
    }
  }
}

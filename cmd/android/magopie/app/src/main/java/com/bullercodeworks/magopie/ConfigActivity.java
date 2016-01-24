package com.bullercodeworks.magopie;

import android.app.Activity;
import android.os.Bundle;
import android.view.View;
import android.widget.Button;
import android.widget.EditText;
import android.widget.TextView;

/**
 * Created by brbuller on 1/23/16.
 */
public class ConfigActivity extends Activity {
  public State state;
  EditText txtServerURL;
  EditText txtApiToken;
  Button btnSaveConfig;
  @Override
  protected void onCreate(Bundle savedInstanceState) {
    super.onCreate(savedInstanceState);
    setContentView(R.layout.content_config);
    state = new State();
    btnSaveConfig = (Button)findViewById(R.id.btnSaveConfig);
    txtServerURL = (EditText)findViewById(R.id.txtServerURL);
    txtApiToken = (EditText)findViewById(R.id.txtApiToken);
    btnSaveConfig.setOnClickListener(new View.OnClickListener() {
      @Override
      public void onClick(View view) {
        state.ServerURL = txtServerURL.getText().toString();
        state.ApiToken = txtApiToken.getText().toString();

        finish();
      }
    });
  }

  @Override
  public void onResume() {
    super.onResume();
    state.load();
    txtServerURL.setText(state.ServerURL);
    txtApiToken.setText(state.ApiToken);
  }

  @Override
  public void onPause() {
    state.save();
    super.onPause();
  }

  @Override
  public void onBackPressed() {
    if("".equals(txtServerURL.getText().toString())) {
      this.finishAffinity();
    } else {
      super.onBackPressed();
    }
  }
}

package com.bullercodeworks.magopie;

import android.animation.Animator;
import android.animation.AnimatorListenerAdapter;
import android.app.Activity;
import android.app.ProgressDialog;
import android.content.Context;
import android.content.Intent;
import android.os.Bundle;
import android.view.KeyEvent;
import android.view.Menu;
import android.view.MenuItem;
import android.view.View;
import android.view.WindowManager;
import android.view.inputmethod.EditorInfo;
import android.view.inputmethod.InputMethodManager;
import android.widget.Button;
import android.widget.EditText;
import android.widget.ImageView;
import android.widget.LinearLayout;
import android.widget.ListView;
import android.widget.TextView;

import com.bullercodeworks.magopie.adapters.ResultAdapter;

import go.magopie.Magopie;

public class MainActivity extends Activity {
  public State state;

  private LinearLayout layoutSearch;
  private ImageView magoLogo;
  private EditText inpSearch;
  private Button btnSearch;
  private Button btnAbout;

  private LinearLayout layoutResults;
  private EditText inpResSearch;
  private ResultAdapter resultAdapter;
  private ListView resultList;

  private int mShortAnimationDuration;

  @Override
  protected void onCreate(Bundle savedInstanceState) {
    super.onCreate(savedInstanceState);
    setContentView(R.layout.content_main);

    state = new State();
    mShortAnimationDuration = getResources().getInteger(android.R.integer.config_shortAnimTime);

    layoutSearch = (LinearLayout)findViewById(R.id.layoutSearch);
    magoLogo = (ImageView)findViewById(R.id.imgMagoLogo);
    inpSearch = (EditText)findViewById(R.id.inpSearch);
    btnSearch = (Button)findViewById(R.id.btnSearch);
    btnAbout = (Button)findViewById(R.id.btnAbout);
    TextView.OnEditorActionListener searchOnDone = new TextView.OnEditorActionListener() {
      @Override
      public boolean onEditorAction(TextView v, int actionId, KeyEvent event) {
        System.out.println("Triggering search...");
        System.out.println(actionId);
        if (actionId == EditorInfo.IME_ACTION_DONE) {
          doSearch();
          return true;
        }
        return false;
      }
    };
    inpSearch.setOnEditorActionListener(searchOnDone);

    btnSearch.setOnClickListener(new View.OnClickListener() {
      @Override
      public void onClick(View view) {
        doSearch();
      }
    });
    btnAbout.setOnClickListener(new View.OnClickListener() {
      @Override
      public void onClick(View view) {
        loadAboutActivity();
      }
    });

    layoutResults = (LinearLayout)findViewById(R.id.layoutResults);
    layoutResults.setVisibility(View.GONE);

    inpResSearch = (EditText)findViewById(R.id.inpResultSearch);
    inpResSearch.setOnEditorActionListener(searchOnDone);

    if(resultAdapter == null) {
      resultAdapter = new ResultAdapter(state, this);
    }
    resultList = (ListView)findViewById(R.id.resultList);
    resultList.setAdapter(resultAdapter);
    resultAdapter.notifyDataSetChanged();
  }

  @Override
  public void onResume() {
    super.onResume();
    state.load();
    if("".equals(state.ServerURL)) {
        loadConfigActivity();
    }
  }

  @Override
  public void onPause() {
    state.save();
    super.onPause();
  }

  @Override
  public void onBackPressed() {
    if (layoutResults.getVisibility() == View.VISIBLE) {
      switchToSearch();
    } else {
      super.onBackPressed();
    }
  }

  @Override
  public boolean onCreateOptionsMenu(Menu menu) {
    // Inflate the menu; this adds items to the action bar if it is present.
    getMenuInflater().inflate(R.menu.menu_main, menu);
    return true;
  }

  @Override
  public boolean onOptionsItemSelected(MenuItem item) {
    // Handle action bar item clicks here. The action bar will
    // automatically handle clicks on the Home/Up button, so long
    // as you specify a parent activity in AndroidManifest.xml.
    int id = item.getItemId();

    //noinspection SimplifiableIfStatement
    if (id == R.id.action_settings) {
      loadConfigActivity();
      return true;
    }

    return super.onOptionsItemSelected(item);
  }

  public void doSearch() {
    String srchTerm = inpSearch.getText().toString();
    if(layoutResults.getVisibility() == View.VISIBLE) {
      srchTerm = inpResSearch.getText().toString();
    }
    /*
    ProgressDialog progress = new ProgressDialog(this);
    progress.setTitle("Loading");
    progress.setMessage("Wait while loading...");
    progress.show();
    */
    Magopie.TorrentCollection wrk = Magopie.NewClient(state.ServerURL, state.ApiToken).Search(srchTerm);
    state.results.clear();
    for(int i = 0; i < wrk.Length(); i++) {
      state.results.add(wrk.Get(i));
    }
    resultAdapter.notifyDataSetChanged();
    /*
    progress.dismiss();
    */
    if(layoutResults.getVisibility() != View.VISIBLE) {
      switchToResults();
    } else {
      // Reload the results list
      switchToSearch();
      switchToResults();
    }
    InputMethodManager imm = (InputMethodManager)getSystemService(Context.INPUT_METHOD_SERVICE);
    imm.hideSoftInputFromWindow(layoutResults.getWindowToken(), 0);
  }

  public void switchToResults() {
    layoutResults.setAlpha(0f);
    layoutResults.setVisibility(View.VISIBLE);
    layoutResults.animate()
        .alpha(1f)
        .setDuration(mShortAnimationDuration)
        .setListener(null);
    layoutSearch.animate()
        .alpha(0f)
        .setDuration(mShortAnimationDuration)
        .setListener(new AnimatorListenerAdapter() {
          @Override
          public void onAnimationEnd(Animator animation) {
            layoutSearch.setVisibility(View.GONE);
          }
        });
  }

  public void switchToSearch() {
    layoutSearch.setAlpha(0f);
    layoutSearch.setVisibility(View.VISIBLE);
    layoutSearch.animate()
        .alpha(1f)
        .setDuration(mShortAnimationDuration)
        .setListener(null);
    layoutResults.animate()
        .alpha(0f)
        .setDuration(mShortAnimationDuration)
        .setListener(new AnimatorListenerAdapter() {
          @Override
          public void onAnimationEnd(Animator animation) {
            layoutResults.setVisibility(View.GONE);
          }
        });
    if(inpSearch.requestFocus()) {
      getWindow().setSoftInputMode(WindowManager.LayoutParams.SOFT_INPUT_STATE_ALWAYS_VISIBLE);
    }
  }

  public void loadConfigActivity() {
    Intent configIntent = new Intent(this, ConfigActivity.class);
    startActivity(configIntent);
  }

  public void loadAboutActivity() {
    Intent aboutIntent = new Intent(this, AboutActivity.class);
    startActivity(aboutIntent);
  }
}
